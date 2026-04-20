# PrioritiesApi

All URIs are relative to *http://localhost:8080/api*

| Method | HTTP request | Description |
|------------- | ------------- | -------------|
| [**prioritiesGet**](PrioritiesApi.md#prioritiesget) | **GET** /priorities | List priorities |
| [**prioritiesIdDelete**](PrioritiesApi.md#prioritiesiddelete) | **DELETE** /priorities/{id} | Delete priority |
| [**prioritiesIdGet**](PrioritiesApi.md#prioritiesidget) | **GET** /priorities/{id} | Get priority by ID |
| [**prioritiesIdPut**](PrioritiesApi.md#prioritiesidput) | **PUT** /priorities/{id} | Update priority |
| [**prioritiesPost**](PrioritiesApi.md#prioritiespost) | **POST** /priorities | Create status |



## prioritiesGet

> WebListPrioritiesResponse prioritiesGet()

List priorities

Get all priorities in dictionary

### Example

```ts
import {
  Configuration,
  PrioritiesApi,
} from '';
import type { PrioritiesGetRequest } from '';

async function example() {
  console.log("🚀 Testing  SDK...");
  const api = new PrioritiesApi();

  try {
    const data = await api.prioritiesGet();
    console.log(data);
  } catch (error) {
    console.error(error);
  }
}

// Run the test
example().catch(console.error);
```

### Parameters

This endpoint does not need any parameter.

### Return type

[**WebListPrioritiesResponse**](WebListPrioritiesResponse.md)

### Authorization

No authorization required

### HTTP request headers

- **Content-Type**: Not defined
- **Accept**: `application/json`


### HTTP response details
| Status code | Description | Response headers |
|-------------|-------------|------------------|
| **200** | OK |  -  |
| **401** | Unauthorized |  -  |
| **403** | Forbidden |  -  |
| **500** | Internal Server Error |  -  |

[[Back to top]](#) [[Back to API list]](../README.md#api-endpoints) [[Back to Model list]](../README.md#models) [[Back to README]](../README.md)


## prioritiesIdDelete

> prioritiesIdDelete(id)

Delete priority

Delete priority dictionary entry by ID

### Example

```ts
import {
  Configuration,
  PrioritiesApi,
} from '';
import type { PrioritiesIdDeleteRequest } from '';

async function example() {
  console.log("🚀 Testing  SDK...");
  const api = new PrioritiesApi();

  const body = {
    // number | Priority ID
    id: 56,
  } satisfies PrioritiesIdDeleteRequest;

  try {
    const data = await api.prioritiesIdDelete(body);
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
| **id** | `number` | Priority ID | [Defaults to `undefined`] |

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
| **403** | Forbidden |  -  |
| **500** | Internal Server Error |  -  |

[[Back to top]](#) [[Back to API list]](../README.md#api-endpoints) [[Back to Model list]](../README.md#models) [[Back to README]](../README.md)


## prioritiesIdGet

> WebPriority prioritiesIdGet(id)

Get priority by ID

Get priority dictionary entry by ID

### Example

```ts
import {
  Configuration,
  PrioritiesApi,
} from '';
import type { PrioritiesIdGetRequest } from '';

async function example() {
  console.log("🚀 Testing  SDK...");
  const api = new PrioritiesApi();

  const body = {
    // number | Priority ID
    id: 56,
  } satisfies PrioritiesIdGetRequest;

  try {
    const data = await api.prioritiesIdGet(body);
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
| **id** | `number` | Priority ID | [Defaults to `undefined`] |

### Return type

[**WebPriority**](WebPriority.md)

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
| **403** | Forbidden |  -  |
| **500** | Internal Server Error |  -  |

[[Back to top]](#) [[Back to API list]](../README.md#api-endpoints) [[Back to Model list]](../README.md#models) [[Back to README]](../README.md)


## prioritiesIdPut

> WebPriority prioritiesIdPut(id, request)

Update priority

Update priority dictionary entry by ID

### Example

```ts
import {
  Configuration,
  PrioritiesApi,
} from '';
import type { PrioritiesIdPutRequest } from '';

async function example() {
  console.log("🚀 Testing  SDK...");
  const api = new PrioritiesApi();

  const body = {
    // number | Priority ID
    id: 56,
    // WebUpdatePriorityRequest | Priority data
    request: ...,
  } satisfies PrioritiesIdPutRequest;

  try {
    const data = await api.prioritiesIdPut(body);
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
| **id** | `number` | Priority ID | [Defaults to `undefined`] |
| **request** | [WebUpdatePriorityRequest](WebUpdatePriorityRequest.md) | Priority data | |

### Return type

[**WebPriority**](WebPriority.md)

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
| **403** | Forbidden |  -  |
| **500** | Internal Server Error |  -  |

[[Back to top]](#) [[Back to API list]](../README.md#api-endpoints) [[Back to Model list]](../README.md#models) [[Back to README]](../README.md)


## prioritiesPost

> WebPriority prioritiesPost(request)

Create status

Create new status dictionary entry

### Example

```ts
import {
  Configuration,
  PrioritiesApi,
} from '';
import type { PrioritiesPostRequest } from '';

async function example() {
  console.log("🚀 Testing  SDK...");
  const api = new PrioritiesApi();

  const body = {
    // WebCreatePriorityRequest | Priority data
    request: ...,
  } satisfies PrioritiesPostRequest;

  try {
    const data = await api.prioritiesPost(body);
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
| **request** | [WebCreatePriorityRequest](WebCreatePriorityRequest.md) | Priority data | |

### Return type

[**WebPriority**](WebPriority.md)

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
| **403** | Forbidden |  -  |
| **500** | Internal Server Error |  -  |

[[Back to top]](#) [[Back to API list]](../README.md#api-endpoints) [[Back to Model list]](../README.md#models) [[Back to README]](../README.md)

