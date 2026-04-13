# TeamsApi

All URIs are relative to *http://localhost:8080/api*

| Method | HTTP request | Description |
|------------- | ------------- | -------------|
| [**teamsGet**](TeamsApi.md#teamsget) | **GET** /teams | List teams |
| [**teamsIdDelete**](TeamsApi.md#teamsiddelete) | **DELETE** /teams/{id} | Delete team |
| [**teamsIdGet**](TeamsApi.md#teamsidget) | **GET** /teams/{id} | Get team by ID |
| [**teamsIdPut**](TeamsApi.md#teamsidput) | **PUT** /teams/{id} | Update team |
| [**teamsPost**](TeamsApi.md#teamspost) | **POST** /teams | Create team |



## teamsGet

> WebListTeamsResponse teamsGet()

List teams

Get all teams in dictionary

### Example

```ts
import {
  Configuration,
  TeamsApi,
} from '';
import type { TeamsGetRequest } from '';

async function example() {
  console.log("🚀 Testing  SDK...");
  const api = new TeamsApi();

  try {
    const data = await api.teamsGet();
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

[**WebListTeamsResponse**](WebListTeamsResponse.md)

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


## teamsIdDelete

> teamsIdDelete(id)

Delete team

Delete team dictionary entry by ID

### Example

```ts
import {
  Configuration,
  TeamsApi,
} from '';
import type { TeamsIdDeleteRequest } from '';

async function example() {
  console.log("🚀 Testing  SDK...");
  const api = new TeamsApi();

  const body = {
    // number | Team ID
    id: 56,
  } satisfies TeamsIdDeleteRequest;

  try {
    const data = await api.teamsIdDelete(body);
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
| **id** | `number` | Team ID | [Defaults to `undefined`] |

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


## teamsIdGet

> WebTeam teamsIdGet(id)

Get team by ID

Get team dictionary entry by ID

### Example

```ts
import {
  Configuration,
  TeamsApi,
} from '';
import type { TeamsIdGetRequest } from '';

async function example() {
  console.log("🚀 Testing  SDK...");
  const api = new TeamsApi();

  const body = {
    // number | Team ID
    id: 56,
  } satisfies TeamsIdGetRequest;

  try {
    const data = await api.teamsIdGet(body);
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
| **id** | `number` | Team ID | [Defaults to `undefined`] |

### Return type

[**WebTeam**](WebTeam.md)

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


## teamsIdPut

> WebTeam teamsIdPut(id, request)

Update team

Update team dictionary entry by ID

### Example

```ts
import {
  Configuration,
  TeamsApi,
} from '';
import type { TeamsIdPutRequest } from '';

async function example() {
  console.log("🚀 Testing  SDK...");
  const api = new TeamsApi();

  const body = {
    // number | Team ID
    id: 56,
    // WebUpdateTeamRequest | Team data
    request: ...,
  } satisfies TeamsIdPutRequest;

  try {
    const data = await api.teamsIdPut(body);
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
| **id** | `number` | Team ID | [Defaults to `undefined`] |
| **request** | [WebUpdateTeamRequest](WebUpdateTeamRequest.md) | Team data | |

### Return type

[**WebTeam**](WebTeam.md)

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


## teamsPost

> WebTeam teamsPost(request)

Create team

Create new team dictionary entry

### Example

```ts
import {
  Configuration,
  TeamsApi,
} from '';
import type { TeamsPostRequest } from '';

async function example() {
  console.log("🚀 Testing  SDK...");
  const api = new TeamsApi();

  const body = {
    // WebCreateTeamRequest | Team data
    request: ...,
  } satisfies TeamsPostRequest;

  try {
    const data = await api.teamsPost(body);
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
| **request** | [WebCreateTeamRequest](WebCreateTeamRequest.md) | Team data | |

### Return type

[**WebTeam**](WebTeam.md)

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

