package service

import (
	"errors"
	"fmt"
	"regexp"
	"strings"

	mapset "github.com/deckarep/golang-set"
	clusterDataStore "github.com/stackrox/rox/central/cluster/datastore"
	deploymentMappings "github.com/stackrox/rox/central/deployment/index/mappings"
	imageMappings "github.com/stackrox/rox/central/image/index/mappings"
	notifierStore "github.com/stackrox/rox/central/notifier/store"
	"github.com/stackrox/rox/central/searchbasedpolicies/matcher"
	"github.com/stackrox/rox/generated/storage"
	"github.com/stackrox/rox/pkg/errorhelpers"
	"github.com/stackrox/rox/pkg/policies"
)

func newPolicyValidator(notifierStorage notifierStore.Store, clusterStorage clusterDataStore.DataStore) *policyValidator {
	return &policyValidator{
		notifierStorage:      notifierStorage,
		clusterStorage:       clusterStorage,
		nameValidator:        regexp.MustCompile(`^[^\n\r\$]{5,64}$`),
		descriptionValidator: regexp.MustCompile(`^[^\$]{1,256}$`),
	}
}

// policyValidator validates the incoming policy.
type policyValidator struct {
	notifierStorage      notifierStore.Store
	clusterStorage       clusterDataStore.DataStore
	nameValidator        *regexp.Regexp
	descriptionValidator *regexp.Regexp
}

func (s *policyValidator) validate(policy *storage.Policy) error {
	s.removeEnforcementsForMissingLifecycles(policy)

	errorList := errorhelpers.NewErrorList("policy invalid")
	errorList.AddError(s.validateName(policy))
	errorList.AddError(s.validateDescription(policy))
	errorList.AddError(s.validateCompilableForLifecycle(policy))
	errorList.AddError(s.validateSeverity(policy))
	errorList.AddError(s.validateCategories(policy))
	errorList.AddError(s.validateScopes(policy))
	errorList.AddError(s.validateWhitelists(policy))
	errorList.AddError(s.validateCapabilities(policy))
	return errorList.ToError()
}

func (s *policyValidator) validateName(policy *storage.Policy) error {
	if policy.GetName() == "" || !s.nameValidator.MatchString(policy.GetName()) {
		return errors.New("policy must have a name, at least 5 chars long, and contain no punctuation or special characters")
	}
	return nil
}

func (s *policyValidator) validateDescription(policy *storage.Policy) error {
	if policy.GetDescription() != "" && !s.descriptionValidator.MatchString(policy.GetDescription()) {
		return errors.New("description, when present, should be of sentence form, and not contain more than 200 characters")
	}
	return nil
}

func (s *policyValidator) validateCompilableForLifecycle(policy *storage.Policy) error {
	if len(policy.GetLifecycleStages()) == 0 {
		return fmt.Errorf("a policy must apply to at least one lifecycle stage")
	}

	errorList := errorhelpers.NewErrorList("error validating lifecycle stage")
	if policies.AppliesAtBuildTime(policy) {
		errorList.AddError(compilesForBuildTime(policy))
	}
	if policies.AppliesAtDeployTime(policy) {
		errorList.AddError(compilesForDeployTime(policy))
	}
	if policies.AppliesAtRunTime(policy) {
		errorList.AddError(compilesForRunTime(policy))
	}
	return errorList.ToError()
}

func (s *policyValidator) removeEnforcementsForMissingLifecycles(policy *storage.Policy) {
	if !policies.AppliesAtBuildTime(policy) {
		removeEnforcementForLifecycle(policy, storage.LifecycleStage_BUILD)
	}
	if !policies.AppliesAtDeployTime(policy) {
		removeEnforcementForLifecycle(policy, storage.LifecycleStage_DEPLOY)
	}
	if !policies.AppliesAtRunTime(policy) {
		removeEnforcementForLifecycle(policy, storage.LifecycleStage_RUNTIME)
	}
}

func (s *policyValidator) validateSeverity(policy *storage.Policy) error {
	if policy.GetSeverity() == storage.Severity_UNSET_SEVERITY {
		return errors.New("a policy must have a severity")
	}
	return nil
}

func (s *policyValidator) validateCapabilities(policy *storage.Policy) error {
	set := mapset.NewSet()
	for _, s := range policy.GetFields().GetAddCapabilities() {
		set.Add(s)
	}
	var duplicates []string
	for _, s := range policy.GetFields().GetDropCapabilities() {
		if set.Contains(s) {
			duplicates = append(duplicates, s)
		}
	}
	if len(duplicates) != 0 {
		return fmt.Errorf("Capabilities '%s' cannot be included in both add and drop", strings.Join(duplicates, ","))
	}
	return nil
}

func (s *policyValidator) validateCategories(policy *storage.Policy) error {
	if len(policy.GetCategories()) == 0 {
		return errors.New("a policy must have at least one category")
	}
	categorySet := make(map[string]struct{})
	for _, c := range policy.GetCategories() {
		categorySet[c] = struct{}{}
	}
	if len(categorySet) != len(policy.GetCategories()) {
		return errors.New("a policy cannot contain duplicate categories")
	}
	return nil
}

func (s *policyValidator) validateNotifiers(policy *storage.Policy) error {
	for _, n := range policy.GetNotifiers() {
		_, exists, err := s.notifierStorage.GetNotifier(n)
		if err != nil {
			return fmt.Errorf("error checking if notifier %s is valid", n)
		}
		if !exists {
			return fmt.Errorf("notifier %s does not exist", n)
		}
	}
	return nil
}

func (s *policyValidator) validateScopes(policy *storage.Policy) error {
	for _, scope := range policy.GetScope() {
		if err := s.validateScope(scope); err != nil {
			return err
		}
	}
	return nil
}

func (s *policyValidator) validateWhitelists(policy *storage.Policy) error {
	for _, whitelist := range policy.GetWhitelists() {
		if err := s.validateWhitelist(whitelist); err != nil {
			return err
		}
	}
	return nil
}

func (s *policyValidator) validateWhitelist(whitelist *storage.Whitelist) error {
	// TODO(cgorman) once we have real whitelist support in UI, add validation for whitelist name
	if whitelist.GetContainer() == nil && whitelist.GetDeployment() == nil && whitelist.GetImage() == nil {
		return errors.New("all whitelists must have some criteria to match on")
	}
	if whitelist.GetContainer() != nil {
		if err := s.validateContainerWhitelist(whitelist); err != nil {
			return err
		}
	}
	if whitelist.GetDeployment() != nil {
		if err := s.validateDeploymentWhitelist(whitelist); err != nil {
			return err
		}
	}
	if whitelist.GetImage() != nil {
		if whitelist.GetImage().GetName() == "" {
			return fmt.Errorf("image whitelist must have nonempty name")
		}
	}
	return nil
}

func (s *policyValidator) validateContainerWhitelist(whitelist *storage.Whitelist) error {
	imageName := whitelist.GetContainer().GetImageName()
	if imageName == nil {
		return errors.New("if container whitelist is defined, then image name must also be defined")
	}
	if imageName.GetRegistry() == "" && imageName.GetRemote() == "" && imageName.GetTag() == "" {
		return errors.New("at least one field of image name must be populated (registry, remote, tag)")
	}
	return nil
}

func (s *policyValidator) validateDeploymentWhitelist(whitelist *storage.Whitelist) error {
	deployment := whitelist.GetDeployment()
	if deployment.GetScope() == nil && deployment.GetName() == "" {
		return errors.New("at least one field of deployment whitelist must be defined")
	}
	if deployment.GetScope() != nil {
		if err := s.validateScope(deployment.GetScope()); err != nil {
			return err
		}
	}
	return nil
}

func (s *policyValidator) validateScope(scope *storage.Scope) error {
	if scope.GetCluster() == "" && scope.GetNamespace() == "" && scope.GetLabel() == nil {
		return errors.New("scope must have at least one field populated")
	}
	return nil
}

func compilesForBuildTime(policy *storage.Policy) error {
	m, err := matcher.ForPolicy(policy, imageMappings.OptionsMap, nil)
	if err != nil {
		return fmt.Errorf("policy configuration is invalid for build time: %s", err)
	}
	if m == nil {
		return errors.New("build time policy contains no image constraints")
	}
	return nil
}

func compilesForDeployTime(policy *storage.Policy) error {
	m, err := matcher.ForPolicy(policy, deploymentMappings.OptionsMap, nil)
	if err != nil {
		return fmt.Errorf("policy configuration is invalid for deploy time: %s", err)
	}
	if m == nil {
		return fmt.Errorf("deploy time policy contains no constraints")
	}
	if policy.GetFields().GetProcessPolicy() != nil {
		return errors.New("deploy time policy cannot contain runtime fields")
	}
	return nil
}

func compilesForRunTime(policy *storage.Policy) error {
	m, err := matcher.ForPolicy(policy, deploymentMappings.OptionsMap, nil)
	if err != nil {
		return fmt.Errorf("policy configuration is invalid for run time: %s", err)
	}
	if m == nil {
		return errors.New("run time policy contains no constraints")
	}
	if policy.GetFields().GetProcessPolicy() == nil {
		return errors.New("run time policy must contain runtime specific constraints")
	}
	return nil
}

var enforcementToLifecycle = map[storage.EnforcementAction]storage.LifecycleStage{
	storage.EnforcementAction_FAIL_BUILD_ENFORCEMENT:                    storage.LifecycleStage_BUILD,
	storage.EnforcementAction_SCALE_TO_ZERO_ENFORCEMENT:                 storage.LifecycleStage_DEPLOY,
	storage.EnforcementAction_UNSATISFIABLE_NODE_CONSTRAINT_ENFORCEMENT: storage.LifecycleStage_DEPLOY,
	storage.EnforcementAction_KILL_POD_ENFORCEMENT:                      storage.LifecycleStage_RUNTIME,
}

func removeEnforcementForLifecycle(policy *storage.Policy, stage storage.LifecycleStage) {
	newActions := policy.EnforcementActions[:0]
	for _, ea := range policy.GetEnforcementActions() {
		if enforcementToLifecycle[ea] != stage {
			newActions = append(newActions, ea)
		}
	}
	policy.EnforcementActions = newActions
}
