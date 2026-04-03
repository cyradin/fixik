# IncidentsApi

All URIs are relative to *http://localhost:8080/api*

| Method | HTTP request | Description |
|------------- | ------------- | -------------|
| [**incidentsGet**](IncidentsApi.md#incidentsget) | **GET** /incidents | List incidents |
| [**incidentsIdCommentsGet**](IncidentsApi.md#incidentsidcommentsget) | **GET** /incidents/{id}/comments | Get incident comments |
| [**incidentsIdCommentsPost**](IncidentsApi.md#incidentsidcommentspost) | **POST** /incidents/{id}/comments | Create incident comment |
| [**incidentsIdDelete**](IncidentsApi.md#incidentsiddelete) | **DELETE** /incidents/{id} | Delete incident |
| [**incidentsIdGet**](IncidentsApi.md#incidentsidget) | **GET** /incidents/{id} | Get incident |
| [**incidentsIdPatch**](IncidentsApi.md#incidentsidpatch) | **PATCH** /incidents/{id} | Update incident |
| [**incidentsPost**](IncidentsApi.md#incidentspost) | **POST** /incidents | Create incident |



## incidentsGet

> WebIncidentListResponse incidentsGet(limit, offset)

List incidents

Get all incidents with pagination

### Example

```ts
import {
  Configuration,
  IncidentsApi,
} from '';
import type { IncidentsGetRequest } from '';

async function example() {
  console.log("🚀 Testing  SDK...");
  const api = new IncidentsApi();

  const body = {
    // number | Limit (optional)
    limit: 56,
    // number | Offset (optional)
    offset: 56,
  } satisfies IncidentsGetRequest;

  try {
    const data = await api.incidentsGet(body);
    console.log(data);
  } catch (error) {
    console.error(error);
  }
}

// Run the test
example().catch(console.error);
```

### Parameters


| Name | Type | Description  | Notes |
|------------- | ------------- | ------------- | -------------|
| **limit** | `number` | Limit | [Optional] [Defaults to `100`] |
| **offset** | `number` | Offset | [Optional] [Defaults to `0`] |

### Return type

[**WebIncidentListResponse**](WebIncidentListResponse.md)

### Authorization

No authorization required

### HTTP request headers

- **Content-Type**: Not defined
- **Accept**: `application/json`


### HTTP response details
| Status code | Description | Response headers |
|-------------|-------------|------------------|
| **200** | OK |  -  |
| **400** | Bad Request |  -  |
| **401** | Unauthorized |  -  |
| **500** | Internal Server Error |  -  |

[[Back to top]](#) [[Back to API list]](../README.md#api-endpoints) [[Back to Model list]](../README.md#models) [[Back to README]](../README.md)


## incidentsIdCommentsGet

> WebIncidentCommentListResponse incidentsIdCommentsGet(id, limit, offset)

Get incident comments

Get all comments for an incident with pagination

### Example

```ts
import {
  Configuration,
  IncidentsApi,
} from '';
import type { IncidentsIdCommentsGetRequest } from '';

async function example() {
  console.log("🚀 Testing  SDK...");
  const api = new IncidentsApi();

  const body = {
    // number | Incident ID
    id: 56,
    // number | Limit (optional)
    limit: 56,
    // number | Offset (optional)
    offset: 56,
  } satisfies IncidentsIdCommentsGetRequest;

  try {
    const data = await api.incidentsIdCommentsGet(body);
    console.log(data);
  } catch (error) {
    console.error(error);
  }
}

// Run the test
example().catch(console.error);
```

### Parameters


| Name | Type | Description  | Notes |
|------------- | ------------- | ------------- | -------------|
| **id** | `number` | Incident ID | [Defaults to `undefined`] |
| **limit** | `number` | Limit | [Optional] [Defaults to `100`] |
| **offset** | `number` | Offset | [Optional] [Defaults to `0`] |

### Return type

[**WebIncidentCommentListResponse**](WebIncidentCommentListResponse.md)

### Authorization

No authorization required

### HTTP request headers

- **Content-Type**: Not defined
- **Accept**: `application/json`


### HTTP response details
| Status code | Description | Response headers |
|-------------|-------------|------------------|
| **200** | OK |  -  |
| **400** | Bad Request |  -  |
| **401** | Unauthorized |  -  |
| **500** | Internal Server Error |  -  |

[[Back to top]](#) [[Back to API list]](../README.md#api-endpoints) [[Back to Model list]](../README.md#models) [[Back to README]](../README.md)


## incidentsIdCommentsPost

> WebIncidentComment incidentsIdCommentsPost(id, request)

Create incident comment

Create a new comment for an incident

### Example

```ts
import {
  Configuration,
  IncidentsApi,
} from '';
import type { IncidentsIdCommentsPostRequest } from '';

async function example() {
  console.log("🚀 Testing  SDK...");
  const api = new IncidentsApi();

  const body = {
    // number | Incident ID
    id: 56,
    // WebIncidentCommentCreateRequest | Comment data
    request: ...,
  } satisfies IncidentsIdCommentsPostRequest;

  try {
    const data = await api.incidentsIdCommentsPost(body);
    console.log(data);
  } catch (error) {
    console.error(error);
  }
}

// Run the test
example().catch(console.error);
```

### Parameters


| Name | Type | Description  | Notes |
|------------- | ------------- | ------------- | -------------|
| **id** | `number` | Incident ID | [Defaults to `undefined`] |
| **request** | [WebIncidentCommentCreateRequest](WebIncidentCommentCreateRequest.md) | Comment data | |

### Return type

[**WebIncidentComment**](WebIncidentComment.md)

### Authorization

No authorization required

### HTTP request headers

- **Content-Type**: `application/json`
- **Accept**: `application/json`


### HTTP response details
| Status code | Description | Response headers |
|-------------|-------------|------------------|
| **200** | OK |  -  |
| **400** | Bad Request |  -  |
| **401** | Unauthorized |  -  |
| **500** | Internal Server Error |  -  |

[[Back to top]](#) [[Back to API list]](../README.md#api-endpoints) [[Back to Model list]](../README.md#models) [[Back to README]](../README.md)


## incidentsIdDelete

> incidentsIdDelete(id)

Delete incident

Delete incident

### Example

```ts
import {
  Configuration,
  IncidentsApi,
} from '';
import type { IncidentsIdDeleteRequest } from '';

async function example() {
  console.log("🚀 Testing  SDK...");
  const api = new IncidentsApi();

  const body = {
    // number | Incident ID
    id: 56,
  } satisfies IncidentsIdDeleteRequest;

  try {
    const data = await api.incidentsIdDelete(body);
    console.log(data);
  } catch (error) {
    console.error(error);
  }
}

// Run the test
example().catch(console.error);
```

### Parameters


| Name | Type | Description  | Notes |
|------------- | ------------- | ------------- | -------------|
| **id** | `number` | Incident ID | [Defaults to `undefined`] |

### Return type

`void` (Empty response body)

### Authorization

No authorization required

### HTTP request headers

- **Content-Type**: Not defined
- **Accept**: `application/json`


### HTTP response details
| Status code | Description | Response headers |
|-------------|-------------|------------------|
| **200** | OK |  -  |
| **400** | Bad Request |  -  |
| **401** | Unauthorized |  -  |
| **500** | Internal Server Error |  -  |

[[Back to top]](#) [[Back to API list]](../README.md#api-endpoints) [[Back to Model list]](../README.md#models) [[Back to README]](../README.md)


## incidentsIdGet

> WebIncidentResponse incidentsIdGet(id)

Get incident

Get incident by ID

### Example

```ts
import {
  Configuration,
  IncidentsApi,
} from '';
import type { IncidentsIdGetRequest } from '';

async function example() {
  console.log("🚀 Testing  SDK...");
  const api = new IncidentsApi();

  const body = {
    // number | Incident ID
    id: 56,
  } satisfies IncidentsIdGetRequest;

  try {
    const data = await api.incidentsIdGet(body);
    console.log(data);
  } catch (error) {
    console.error(error);
  }
}

// Run the test
example().catch(console.error);
```

### Parameters


| Name | Type | Description  | Notes |
|------------- | ------------- | ------------- | -------------|
| **id** | `number` | Incident ID | [Defaults to `undefined`] |

### Return type

[**WebIncidentResponse**](WebIncidentResponse.md)

### Authorization

No authorization required

### HTTP request headers

- **Content-Type**: Not defined
- **Accept**: `application/json`


### HTTP response details
| Status code | Description | Response headers |
|-------------|-------------|------------------|
| **200** | OK |  -  |
| **400** | Bad Request |  -  |
| **401** | Unauthorized |  -  |
| **500** | Internal Server Error |  -  |

[[Back to top]](#) [[Back to API list]](../README.md#api-endpoints) [[Back to Model list]](../README.md#models) [[Back to README]](../README.md)


## incidentsIdPatch

> WebIncidentResponse incidentsIdPatch(id, request)

Update incident

Update incident

### Example

```ts
import {
  Configuration,
  IncidentsApi,
} from '';
import type { IncidentsIdPatchRequest } from '';

async function example() {
  console.log("🚀 Testing  SDK...");
  const api = new IncidentsApi();

  const body = {
    // number | Incident ID
    id: 56,
    // WebUpdateIncidentRequest | Incident data
    request: ...,
  } satisfies IncidentsIdPatchRequest;

  try {
    const data = await api.incidentsIdPatch(body);
    console.log(data);
  } catch (error) {
    console.error(error);
  }
}

// Run the test
example().catch(console.error);
```

### Parameters


| Name | Type | Description  | Notes |
|------------- | ------------- | ------------- | -------------|
| **id** | `number` | Incident ID | [Defaults to `undefined`] |
| **request** | [WebUpdateIncidentRequest](WebUpdateIncidentRequest.md) | Incident data | |

### Return type

[**WebIncidentResponse**](WebIncidentResponse.md)

### Authorization

No authorization required

### HTTP request headers

- **Content-Type**: `application/json`
- **Accept**: `application/json`


### HTTP response details
| Status code | Description | Response headers |
|-------------|-------------|------------------|
| **200** | OK |  -  |
| **400** | Bad Request |  -  |
| **401** | Unauthorized |  -  |
| **500** | Internal Server Error |  -  |

[[Back to top]](#) [[Back to API list]](../README.md#api-endpoints) [[Back to Model list]](../README.md#models) [[Back to README]](../README.md)


## incidentsPost

> WebIncidentResponse incidentsPost(request)

Create incident

Create new incident

### Example

```ts
import {
  Configuration,
  IncidentsApi,
} from '';
import type { IncidentsPostRequest } from '';

async function example() {
  console.log("🚀 Testing  SDK...");
  const api = new IncidentsApi();

  const body = {
    // WebCreateIncidentRequest | Incident data
    request: ...,
  } satisfies IncidentsPostRequest;

  try {
    const data = await api.incidentsPost(body);
    console.log(data);
  } catch (error) {
    console.error(error);
  }
}

// Run the test
example().catch(console.error);
```

### Parameters


| Name | Type | Description  | Notes |
|------------- | ------------- | ------------- | -------------|
| **request** | [WebCreateIncidentRequest](WebCreateIncidentRequest.md) | Incident data | |

### Return type

[**WebIncidentResponse**](WebIncidentResponse.md)

### Authorization

No authorization required

### HTTP request headers

- **Content-Type**: `application/json`
- **Accept**: `application/json`


### HTTP response details
| Status code | Description | Response headers |
|-------------|-------------|------------------|
| **200** | OK |  -  |
| **400** | Bad Request |  -  |
| **401** | Unauthorized |  -  |
| **500** | Internal Server Error |  -  |

[[Back to top]](#) [[Back to API list]](../README.md#api-endpoints) [[Back to Model list]](../README.md#models) [[Back to README]](../README.md)

