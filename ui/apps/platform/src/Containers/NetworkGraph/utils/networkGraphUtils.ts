import { EdgeModel, Model, NodeModel } from '@patternfly/react-topology';

import { ListenPort } from 'types/networkFlow.proto';

/* node helper functions */

function getExternalNodeIds(nodes: NodeModel[]): string[] {
    const externalNodeIds =
        nodes?.reduce((acc, curr) => {
            if (curr.data.type === 'INTERNET' || curr.data.type === 'EXTERNAL_SOURCE') {
                return [...acc, curr.id];
            }
            return acc;
        }, [] as string[]) || [];
    return externalNodeIds;
}

export function getNodeById(model: Model, nodeId: string | undefined): NodeModel | undefined {
    return model.nodes?.find((node) => node.id === nodeId);
}

/* edge helper functions */

export function getNumInternalFlows(
    nodes: NodeModel[],
    edges: EdgeModel[],
    deploymentId: string
): number {
    const externalNodeIds = getExternalNodeIds(nodes);
    const numInternalFlows =
        edges?.reduce((acc, edge) => {
            if (
                (edge.source === deploymentId && !externalNodeIds.includes(edge.target || '')) ||
                (edge.target === deploymentId && !externalNodeIds.includes(edge.source || ''))
            ) {
                return acc + 1;
            }
            return acc;
        }, 0) || 0;
    return numInternalFlows;
}

export function getNumExternalFlows(
    nodes: NodeModel[],
    edges: EdgeModel[],
    deploymentId: string
): number {
    const externalNodeIds = getExternalNodeIds(nodes);
    const numExternalFlows =
        edges?.reduce((acc, edge) => {
            if (
                (edge.source === deploymentId && externalNodeIds.includes(edge.target || '')) ||
                (edge.target === deploymentId && externalNodeIds.includes(edge.source || ''))
            ) {
                return acc + 1;
            }
            return acc;
        }, 0) || 0;
    return numExternalFlows;
}

/* deployment helper functions */

export function getListenPorts(nodes: NodeModel[], deploymentId: string): ListenPort[] {
    const deployment = nodes?.find((node) => {
        return node.id === deploymentId;
    });
    if (!deployment) {
        return [];
    }
    return deployment.data.deployment.listenPorts as ListenPort[];
}
