# StatusesApi

All URIs are relative to *http://localhost:8080/api*

| Method | HTTP request | Description |
|------------- | ------------- | -------------|
| [**statusesGet**](StatusesApi.md#statusesget) | **GET** /statuses | List statuses |
| [**statusesIdDelete**](StatusesApi.md#statusesiddelete) | **DELETE** /statuses/{id} | Delete status |
| [**statusesIdGet**](StatusesApi.md#statusesidget) | **GET** /statuses/{id} | Get status by ID |
| [**statusesIdPut**](StatusesApi.md#statusesidput) | **PUT** /statuses/{id} | Update status |
| [**statusesPost**](StatusesApi.md#statusespost) | **POST** /statuses | Create status |



## statusesGet

> WebListStatusesResponse statusesGet()

List statuses

Get all statuses in dictionary

### Example

```ts
import {
  Configuration,
  StatusesApi,
} from '';
import type { StatusesGetRequest } from '';

async function example() {
  console.log("🚀 Testing  SDK...");
  const api = new StatusesApi();

  try {
    const data = await api.statusesGet();
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

[**WebListStatusesResponse**](WebListStatusesResponse.md)

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
| **500** | Internal Server Error |  -  |

[[Back to top]](#) [[Back to API list]](../README.md#api-endpoints) [[Back to Model list]](../README.md#models) [[Back to README]](../README.md)


## statusesIdDelete

> statusesIdDelete(id)

Delete status

Delete status dictionary entry by ID

### Example

```ts
import {
  Configuration,
  StatusesApi,
} from '';
import type { StatusesIdDeleteRequest } from '';

async function example() {
  console.log("🚀 Testing  SDK...");
  const api = new StatusesApi();

  const body = {
    // number | Status ID
    id: 56,
  } satisfies StatusesIdDeleteRequest;

  try {
    const data = await api.statusesIdDelete(body);
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
| **id** | `number` | Status ID | [Defaults to `undefined`] |

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


## statusesIdGet

> WebStatus statusesIdGet(id)

Get status by ID

Get status dictionary entry by ID

### Example

```ts
import {
  Configuration,
  StatusesApi,
} from '';
import type { StatusesIdGetRequest } from '';

async function example() {
  console.log("🚀 Testing  SDK...");
  const api = new StatusesApi();

  const body = {
    // number | Status ID
    id: 56,
  } satisfies StatusesIdGetRequest;

  try {
    const data = await api.statusesIdGet(body);
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
| **id** | `number` | Status ID | [Defaults to `undefined`] |

### Return type

[**WebStatus**](WebStatus.md)

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


## statusesIdPut

> WebStatus statusesIdPut(id, request)

Update status

Update status dictionary entry by ID

### Example

```ts
import {
  Configuration,
  StatusesApi,
} from '';
import type { StatusesIdPutRequest } from '';

async function example() {
  console.log("🚀 Testing  SDK...");
  const api = new StatusesApi();

  const body = {
    // number | Status ID
    id: 56,
    // WebUpdateStatusRequest | Status data
    request: ...,
  } satisfies StatusesIdPutRequest;

  try {
    const data = await api.statusesIdPut(body);
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
| **id** | `number` | Status ID | [Defaults to `undefined`] |
| **request** | [WebUpdateStatusRequest](WebUpdateStatusRequest.md) | Status data | |

### Return type

[**WebStatus**](WebStatus.md)

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


## statusesPost

> WebStatus statusesPost(request)

Create status

Create new status dictionary entry

### Example

```ts
import {
  Configuration,
  StatusesApi,
} from '';
import type { StatusesPostRequest } from '';

async function example() {
  console.log("🚀 Testing  SDK...");
  const api = new StatusesApi();

  const body = {
    // WebCreateStatusRequest | Status data
    request: ...,
  } satisfies StatusesPostRequest;

  try {
    const data = await api.statusesPost(body);
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
| **request** | [WebCreateStatusRequest](WebCreateStatusRequest.md) | Status data | |

### Return type

[**WebStatus**](WebStatus.md)

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

