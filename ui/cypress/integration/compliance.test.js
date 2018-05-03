import { selectors, url as complianceUrl } from './pages/CompliancePage';
import * as api from './apiEndpoints';

describe('Compliance page', () => {
    const setupSingleClusterFixtures = () => {
        cy.server();
        cy.fixture('clusters/single.json').as('singleCluster');
        cy.route('GET', api.clusters.list, '@singleCluster').as('clusters');
        cy.fixture('benchmarks/configs.json').as('configs');
        cy.route('GET', api.benchmarks.configs, '@configs').as('benchConfigs');
        cy.fixture('benchmarks/dockerBenchScans.json').as('dockerBenchScans');
        cy.route('GET', api.benchmarks.cisDockerScans, '@dockerBenchScans').as('scanMetadata');
        cy.fixture('benchmarks/dockerBenchScan1.json').as('dockerBenchScan1');
        cy.route('GET', api.benchmarks.scans, '@dockerBenchScan1').as('benchScan');

        cy.visit(complianceUrl);
        cy.wait(['@clusters', '@benchConfigs', '@scanMetadata', '@benchScan']);
    };

    it('should have selected item in nav bar', () => {
        cy.visit('/');
        cy.get(selectors.compliance).click();
        cy.get(selectors.navLink).click();
        cy.get(selectors.compliance).should('have.class', 'bg-primary-600');
        // first tab selected by default
        cy
            .get(selectors.benchmarkTabs)
            .first()
            .should('have.class', 'tab-active');
    });

    it('should allow scanning initiation', () => {
        cy.server();
        cy.route('POST', api.benchmarks.triggers, {}).as('trigger');
        cy.visit(complianceUrl);
        cy.get(selectors.scanNowButton).as('scanNow');

        cy.get('@scanNow').should('contain', 'Scan now');
        cy.get('@scanNow').click();
        cy.wait('@trigger');
        cy.get('@scanNow').should('not.contain', 'Scan now'); // spinner
    });

    it('should allow to set schedule', () => {
        cy.visit(complianceUrl);

        cy.get('select:first').select('Friday', { force: true });
        cy.get('select:last').select('05:00 PM', { force: true });
        cy.reload(); // retrieve data from the server
        cy.get('select:first').should('have.value', 'Friday');
        cy.get('select:last').should('have.value', '05:00 PM');

        // update schedule
        cy.get('select:last').select('06:00 PM', { force: true });
        cy.reload();
        cy.get('select:last').should('have.value', '06:00 PM');

        // remove schedule
        cy.get('select:first').select('None', { force: true });
        cy.get('select:last').should('have.value', null);
    });

    it('should show scan results', () => {
        setupSingleClusterFixtures();
        cy
            .get(selectors.benchmarkTabs)
            .first()
            .should('contain', 'CIS Docker v1.1.0 Benchmark');
        cy.get(selectors.checkRows).should('have.length', 5);
        cy
            .get(selectors.passColumns)
            .last()
            .should('have.text', '1');
    });

    it('should show benchmark host results', () => {
        setupSingleClusterFixtures();
        cy
            .get(selectors.passColumns)
            .last()
            .click();
        cy
            .get(selectors.hostColumns)
            .should('have.length', 1)
            .and('contain', 'linuxkit-025000000001');
    });
});
