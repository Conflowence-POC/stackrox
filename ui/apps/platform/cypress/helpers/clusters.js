import * as api from '../constants/apiEndpoints';
import { clustersUrl, selectors } from '../constants/ClustersPage';

import { visitFromLeftNavExpandable } from './nav';
import { interactAndWaitForResponses } from './request';
import { visit } from './visit';

const routeMatcherMap = {
    'sensorupgrades/config': {
        method: 'GET',
        url: api.clusters.sensorUpgradesConfig,
    },
    clusters: {
        method: 'GET',
        url: api.clusters.list,
    },
    'cluster-defaults': {
        method: 'GET',
        url: api.clusters.clusterDefaults,
    },
};

// Navigation

/*
 * Reach clusters by interaction from another container.
 * For example, click View All button from System Health.
 */
export function reachClusters(interactionCallback, staticResponseMap) {
    interactAndWaitForResponses(interactionCallback, routeMatcherMap, staticResponseMap);

    cy.get(selectors.clustersListHeading).contains('Clusters');
}

export function visitClustersFromLeftNav() {
    visitFromLeftNavExpandable('Platform Configuration', 'Clusters', routeMatcherMap);

    cy.location('pathname').should('eq', clustersUrl);
    cy.get(selectors.clustersListHeading).contains('Clusters');
}

export function visitClusters(staticResponseMap) {
    visit(clustersUrl, routeMatcherMap, staticResponseMap);

    cy.get(selectors.clustersListHeading).contains('Clusters');
}

export function visitClustersWithFixture(fixturePath) {
    visitClusters({
        clusters: { fixture: fixturePath },
    });
}

export function visitClusterById(clusterId, staticResponseMap) {
    const routeMatcherMapClusterById = {
        'cluster-defaults': {
            method: 'GET',
            url: api.clusters.clusterDefaults,
        },
        cluster: {
            method: 'GET',
            url: `${api.clusters.list}/${clusterId}`,
        },
    };
    visit(`${clustersUrl}/${clusterId}`, routeMatcherMapClusterById, staticResponseMap);

    cy.get(selectors.clustersListHeading).contains('Clusters');
}

export function visitClustersWithFixtureMetadataDatetime(fixturePath, metadata, datetimeISOString) {
    cy.intercept('GET', api.metadata, {
        body: metadata,
    }).as('metadata');

    // For comparison to `lastContact` and `sensorCertExpiry` in clusters fixture.
    const currentDatetime = new Date(datetimeISOString);
    cy.clock(currentDatetime.getTime(), ['Date', 'setInterval']);

    visitClustersWithFixture(fixturePath);

    cy.wait('@metadata');
}

export function visitClusterByNameWithFixture(clusterName, fixturePath) {
    cy.fixture(fixturePath).then(({ clusters, clusterIdToRetentionInfo }) => {
        const cluster = clusters.find(({ name }) => name === clusterName);
        const clusterRetentionInfo = clusterIdToRetentionInfo[cluster.id] ?? null;

        visitClusterById(cluster.id, {
            cluster: { body: { cluster, clusterRetentionInfo } },
        });

        cy.get(selectors.clusterSidePanelHeading).contains(clusterName);
    });
}

export function visitClusterByNameWithFixtureMetadataDatetime(
    clusterName,
    fixturePath,
    metadata,
    datetimeISOString
) {
    cy.fixture(fixturePath).then(({ clusters, clusterIdToRetentionInfo }) => {
        cy.intercept('GET', api.metadata, {
            body: metadata,
        }).as('metadata');

        const cluster = clusters.find(({ name }) => name === clusterName);
        const clusterRetentionInfo = clusterIdToRetentionInfo[cluster.id] ?? null;

        // For comparison to `lastContact` and `sensorCertExpiry` in clusters fixture.
        const currentDatetime = new Date(datetimeISOString);
        cy.clock(currentDatetime.getTime(), ['Date', 'setInterval']);

        visitClusterById(cluster.id, {
            cluster: { body: { cluster, clusterRetentionInfo } },
        });

        cy.wait(['@metadata']);
        cy.get(selectors.clusterSidePanelHeading).contains(clusterName);
    });
}

export function visitDashboardWithNoClusters() {
    cy.intercept('POST', api.graphql('summary_counts'), {
        body: {
            data: {
                clusterCount: 0,
                nodeCount: 3,
                violationCount: 20,
                deploymentCount: 35,
                imageCount: 31,
                secretCount: 15,
            },
        },
    }).as('summary_counts');
    cy.intercept('GET', api.clusters.list, {
        clusters: [],
    }).as('clusters');

    // visitMainDashboard(); // with a count of 0 clusters, app should redirect to the clusters pages
    cy.visit('/main/dashboard'); // with a count of 0 clusters, app should redirect to the clusters pages

    cy.wait(['@summary_counts', '@clusters']);
}
