import React from 'react';
import {
    Alert,
    AlertVariant,
    Badge,
    Bullseye,
    Flex,
    FlexItem,
    Spinner,
    Stack,
    StackItem,
    Tab,
    TabContent,
    Tabs,
    TabTitleText,
    Text,
    TextContent,
    TextVariants,
} from '@patternfly/react-core';
import { EdgeModel, NodeModel } from '@patternfly/react-topology';

import useTabs from 'hooks/patternfly/useTabs';
import useFetchDeployment from 'hooks/useFetchDeployment';
import {
    getListenPorts,
    getNumExternalFlows,
    getNumInternalFlows,
} from '../utils/networkGraphUtils';

import DeploymentDetails from './DeploymentDetails';
import DeploymentNetworkPolicies from './DeploymentNetworkPolicies';
import DeploymentFlows from './DeploymentFlows';

type DeploymentSideBarProps = {
    deploymentId: string;
    nodes: NodeModel[];
    edges: EdgeModel[];
};

function DeploymentSideBar({ deploymentId, nodes, edges }: DeploymentSideBarProps) {
    // component state
    const { deployment, isLoading, error } = useFetchDeployment(deploymentId);
    const { activeKeyTab, onSelectTab } = useTabs({
        defaultTab: 'Details',
    });

    // derived values
    const numExternalFlows = getNumExternalFlows(nodes, edges, deploymentId);
    const numInternalFlows = getNumInternalFlows(nodes, edges, deploymentId);
    const listenPorts = getListenPorts(nodes, deploymentId);

    if (isLoading) {
        return (
            <Bullseye>
                <Spinner isSVG size="lg" />
            </Bullseye>
        );
    }

    if (error) {
        return (
            <Alert isInline variant={AlertVariant.danger} title={error} className="pf-u-mb-lg" />
        );
    }

    return (
        <Stack>
            <StackItem>
                <Flex direction={{ default: 'row' }} className="pf-u-p-md pf-u-mb-0">
                    <FlexItem>
                        <Badge style={{ backgroundColor: 'rgb(0,102,205)' }}>D</Badge>
                    </FlexItem>
                    <FlexItem>
                        <TextContent>
                            <Text component={TextVariants.h1} className="pf-u-font-size-xl">
                                {deployment?.name}
                            </Text>
                        </TextContent>
                        <TextContent>
                            <Text
                                component={TextVariants.h2}
                                className="pf-u-font-size-sm pf-u-color-200"
                            >
                                in &quot;{deployment?.clusterName} / {deployment?.namespace}&quot;
                            </Text>
                        </TextContent>
                    </FlexItem>
                </Flex>
            </StackItem>
            <StackItem>
                <Tabs activeKey={activeKeyTab} onSelect={onSelectTab}>
                    <Tab
                        eventKey="Details"
                        tabContentId="Details"
                        title={<TabTitleText>Details</TabTitleText>}
                    />
                    <Tab
                        eventKey="Flows"
                        tabContentId="Flows"
                        title={<TabTitleText>Flows</TabTitleText>}
                    />
                    <Tab
                        eventKey="Baselines"
                        tabContentId="Baselines"
                        title={<TabTitleText>Baselines</TabTitleText>}
                    />
                    <Tab
                        eventKey="Network policies"
                        tabContentId="Network policies"
                        title={<TabTitleText>Network policies</TabTitleText>}
                    />
                </Tabs>
            </StackItem>
            <StackItem isFilled style={{ overflow: 'auto' }}>
                <TabContent eventKey="Details" id="Details" hidden={activeKeyTab !== 'Details'}>
                    {deployment && (
                        <DeploymentDetails
                            deployment={deployment}
                            numExternalFlows={numExternalFlows}
                            numInternalFlows={numInternalFlows}
                            listenPorts={listenPorts}
                        />
                    )}
                </TabContent>
                <TabContent eventKey="Flows" id="Flows" hidden={activeKeyTab !== 'Flows'}>
                    <DeploymentFlows />
                </TabContent>
                <TabContent
                    eventKey="Baselines"
                    id="Baselines"
                    hidden={activeKeyTab !== 'Baselines'}
                >
                    <div className="pf-u-h-100 pf-u-p-md">TODO: Add Baselines</div>
                </TabContent>
                <TabContent
                    eventKey="Network policies"
                    id="Network policies"
                    hidden={activeKeyTab !== 'Network policies'}
                >
                    <DeploymentNetworkPolicies />
                </TabContent>
            </StackItem>
        </Stack>
    );
}

export default DeploymentSideBar;
