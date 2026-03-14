# ImpactsApi

All URIs are relative to *http://localhost:8080/api*

| Method | HTTP request | Description |
|------------- | ------------- | -------------|
| [**impactsGet**](ImpactsApi.md#impactsget) | **GET** /impacts | List impacts |
| [**impactsIdDelete**](ImpactsApi.md#impactsiddelete) | **DELETE** /impacts/{id} | Delete impact |
| [**impactsIdGet**](ImpactsApi.md#impactsidget) | **GET** /impacts/{id} | Get impact by ID |
| [**impactsIdPut**](ImpactsApi.md#impactsidput) | **PUT** /impacts/{id} | Update impact |
| [**impactsPost**](ImpactsApi.md#impactspost) | **POST** /impacts | Create impact |



## impactsGet

> WebListDictEntitiesResponse impactsGet()

List impacts

Get all impacts in dictionary

### Example

```ts
import {
  Configuration,
  ImpactsApi,
} from '';
import type { ImpactsGetRequest } from '';

async function example() {
  console.log("🚀 Testing  SDK...");
  const api = new ImpactsApi();

  try {
    const data = await api.impactsGet();
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

[**WebListDictEntitiesResponse**](WebListDictEntitiesResponse.md)

### Authorization

No authorization required

### HTTP request headers

- **Content-Type**: Not defined
- **Accept**: `application/json`


### HTTP response details
| Status code | Description | Response headers |
|-------------|-------------|------------------|
| **200** | OK |  -  |
| **500** | Internal Server Error |  -  |

[[Back to top]](#) [[Back to API list]](../README.md#api-endpoints) [[Back to Model list]](../README.md#models) [[Back to README]](../README.md)


## impactsIdDelete

> impactsIdDelete(id)

Delete impact

Delete impact dictionary entry by ID

### Example

```ts
import {
  Configuration,
  ImpactsApi,
} from '';
import type { ImpactsIdDeleteRequest } from '';

async function example() {
  console.log("🚀 Testing  SDK...");
  const api = new ImpactsApi();

  const body = {
    // number | Impact ID
    id: 56,
  } satisfies ImpactsIdDeleteRequest;

  try {
    const data = await api.impactsIdDelete(body);
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
| **id** | `number` | Impact ID | [Defaults to `undefined`] |

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
| **500** | Internal Server Error |  -  |

[[Back to top]](#) [[Back to API list]](../README.md#api-endpoints) [[Back to Model list]](../README.md#models) [[Back to README]](../README.md)


## impactsIdGet

> WebDictEntity impactsIdGet(id)

Get impact by ID

Get impact dictionary entry by ID

### Example

```ts
import {
  Configuration,
  ImpactsApi,
} from '';
import type { ImpactsIdGetRequest } from '';

async function example() {
  console.log("🚀 Testing  SDK...");
  const api = new ImpactsApi();

  const body = {
    // number | Impact ID
    id: 56,
  } satisfies ImpactsIdGetRequest;

  try {
    const data = await api.impactsIdGet(body);
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
| **id** | `number` | Impact ID | [Defaults to `undefined`] |

### Return type

[**WebDictEntity**](WebDictEntity.md)

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
| **500** | Internal Server Error |  -  |

[[Back to top]](#) [[Back to API list]](../README.md#api-endpoints) [[Back to Model list]](../README.md#models) [[Back to README]](../README.md)


## impactsIdPut

> WebDictEntity impactsIdPut(id, request)

Update impact

Update impact dictionary entry by ID

### Example

```ts
import {
  Configuration,
  ImpactsApi,
} from '';
import type { ImpactsIdPutRequest } from '';

async function example() {
  console.log("🚀 Testing  SDK...");
  const api = new ImpactsApi();

  const body = {
    // number | Impact ID
    id: 56,
    // WebUpdateDictEntityRequest | Impact data
    request: ...,
  } satisfies ImpactsIdPutRequest;

  try {
    const data = await api.impactsIdPut(body);
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
| **id** | `number` | Impact ID | [Defaults to `undefined`] |
| **request** | [WebUpdateDictEntityRequest](WebUpdateDictEntityRequest.md) | Impact data | |

### Return type

[**WebDictEntity**](WebDictEntity.md)

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
| **500** | Internal Server Error |  -  |

[[Back to top]](#) [[Back to API list]](../README.md#api-endpoints) [[Back to Model list]](../README.md#models) [[Back to README]](../README.md)


## impactsPost

> WebDictEntity impactsPost(request)

Create impact

Create new impact dictionary entry

### Example

```ts
import {
  Configuration,
  ImpactsApi,
} from '';
import type { ImpactsPostRequest } from '';

async function example() {
  console.log("🚀 Testing  SDK...");
  const api = new ImpactsApi();

  const body = {
    // WebCreateDictEntityRequest | Impact data
    request: ...,
  } satisfies ImpactsPostRequest;

  try {
    const data = await api.impactsPost(body);
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
| **request** | [WebCreateDictEntityRequest](WebCreateDictEntityRequest.md) | Impact data | |

### Return type

[**WebDictEntity**](WebDictEntity.md)

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
| **500** | Internal Server Error |  -  |

[[Back to top]](#) [[Back to API list]](../README.md#api-endpoints) [[Back to Model list]](../README.md#models) [[Back to README]](../README.md)

