import apiBaseUrl from "./api-base-url"
import {
  INodeContainerBaseUrlUpdateRequest,
  INodeCreateRequest,
  INodeUpdateRequest,
  ISettingUpdateRequest,
} from "./api-models"

// Nodes

export async function apiNodesGenerateToken(nodeId: Number) {
  return apiPost(`${apiBaseUrl()}/nodes/${nodeId}/generatetoken`, null)
}

export async function apiNodesCreate(request: INodeCreateRequest) {
  return apiPost(`${apiBaseUrl()}/nodes`, request)
}

export async function apiNodesUpdate(
  nodeId: Number,
  request: INodeUpdateRequest
) {
  return apiPut(`${apiBaseUrl()}/nodes/${nodeId}`, request)
}

export async function apiNodesContainerBaseUrlUpdate(
  nodeId: Number,
  request: INodeContainerBaseUrlUpdateRequest
) {
  return apiPatch(`${apiBaseUrl()}/nodes/${nodeId}`, request)
}

export async function apiNodesDelete(nodeId: Number) {
  return apiDelete(`${apiBaseUrl()}/nodes/${nodeId}`, null)
}

// Settings

export async function apiSettingsUpdate(
  id: string,
  request: ISettingUpdateRequest
) {
  return apiPut(`${apiBaseUrl()}/settings/${id}`, request)
}

// Common

function apiPost(url: string, body: any) {
  return fetch(url, {
    method: "POST",
    ...commonOptions,
    body: body ? JSON.stringify(body) : null,
  })
}

function apiPut(url: string, body: any) {
  return fetch(url, {
    method: "PUT",
    ...commonOptions,
    body: body ? JSON.stringify(body) : null,
  })
}

function apiPatch(url: string, body: any) {
  return fetch(url, {
    method: "PATCH",
    ...commonOptions,
    body: body ? JSON.stringify(body) : null,
  })
}

function apiDelete(url: string, body: any) {
  return fetch(url, {
    method: "DELETE",
    ...commonOptions,
    body: body ? JSON.stringify(body) : null,
  })
}

const commonOptions = {
  headers: { "Content-Type": "application/json" },
}
