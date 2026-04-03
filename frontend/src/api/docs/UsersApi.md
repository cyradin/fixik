# UsersApi

All URIs are relative to *http://localhost:8080/api*

| Method | HTTP request | Description |
|------------- | ------------- | -------------|
| [**usersGet**](UsersApi.md#usersget) | **GET** /users | List users |
| [**usersIdDelete**](UsersApi.md#usersiddelete) | **DELETE** /users/{id} | Delete user |
| [**usersIdGet**](UsersApi.md#usersidget) | **GET** /users/{id} | Get user by ID |
| [**usersIdPasswordPost**](UsersApi.md#usersidpasswordpost) | **POST** /users/{id}/password | Change user password |
| [**usersIdPatch**](UsersApi.md#usersidpatch) | **PATCH** /users/{id} | Update user |
| [**usersPost**](UsersApi.md#userspost) | **POST** /users | Create user |



## usersGet

> WebListUsersResponse usersGet(limit, offset)

List users

Get all users with optional limit and offset

### Example

```ts
import {
  Configuration,
  UsersApi,
} from '';
import type { UsersGetRequest } from '';

async function example() {
  console.log("🚀 Testing  SDK...");
  const api = new UsersApi();

  const body = {
    // number | Limit of items
    limit: 56,
    // number | Offset for pagination
    offset: 56,
  } satisfies UsersGetRequest;

  try {
    const data = await api.usersGet(body);
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
| **limit** | `number` | Limit of items | [Defaults to `100`] |
| **offset** | `number` | Offset for pagination | [Defaults to `0`] |

### Return type

[**WebListUsersResponse**](WebListUsersResponse.md)

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


## usersIdDelete

> usersIdDelete(id)

Delete user

Delete user by ID

### Example

```ts
import {
  Configuration,
  UsersApi,
} from '';
import type { UsersIdDeleteRequest } from '';

async function example() {
  console.log("🚀 Testing  SDK...");
  const api = new UsersApi();

  const body = {
    // number | User ID
    id: 56,
  } satisfies UsersIdDeleteRequest;

  try {
    const data = await api.usersIdDelete(body);
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
| **id** | `number` | User ID | [Defaults to `undefined`] |

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


## usersIdGet

> WebUserResponse usersIdGet(id)

Get user by ID

Get user by ID

### Example

```ts
import {
  Configuration,
  UsersApi,
} from '';
import type { UsersIdGetRequest } from '';

async function example() {
  console.log("🚀 Testing  SDK...");
  const api = new UsersApi();

  const body = {
    // number | User ID
    id: 56,
  } satisfies UsersIdGetRequest;

  try {
    const data = await api.usersIdGet(body);
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
| **id** | `number` | User ID | [Defaults to `undefined`] |

### Return type

[**WebUserResponse**](WebUserResponse.md)

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


## usersIdPasswordPost

> usersIdPasswordPost(id, request)

Change user password

Change current user\&#39;s password

### Example

```ts
import {
  Configuration,
  UsersApi,
} from '';
import type { UsersIdPasswordPostRequest } from '';

async function example() {
  console.log("🚀 Testing  SDK...");
  const api = new UsersApi();

  const body = {
    // number | User ID
    id: 56,
    // WebChangePasswordRequest | Password data
    request: ...,
  } satisfies UsersIdPasswordPostRequest;

  try {
    const data = await api.usersIdPasswordPost(body);
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
| **id** | `number` | User ID | [Defaults to `undefined`] |
| **request** | [WebChangePasswordRequest](WebChangePasswordRequest.md) | Password data | |

### Return type

`void` (Empty response body)

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


## usersIdPatch

> WebUserResponse usersIdPatch(id, request)

Update user

Update user by ID

### Example

```ts
import {
  Configuration,
  UsersApi,
} from '';
import type { UsersIdPatchRequest } from '';

async function example() {
  console.log("🚀 Testing  SDK...");
  const api = new UsersApi();

  const body = {
    // number | User ID
    id: 56,
    // WebUpdateUserRequest | User data
    request: ...,
  } satisfies UsersIdPatchRequest;

  try {
    const data = await api.usersIdPatch(body);
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
| **id** | `number` | User ID | [Defaults to `undefined`] |
| **request** | [WebUpdateUserRequest](WebUpdateUserRequest.md) | User data | |

### Return type

[**WebUserResponse**](WebUserResponse.md)

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


## usersPost

> WebUserResponse usersPost(request)

Create user

Create new user

### Example

```ts
import {
  Configuration,
  UsersApi,
} from '';
import type { UsersPostRequest } from '';

async function example() {
  console.log("🚀 Testing  SDK...");
  const api = new UsersApi();

  const body = {
    // WebCreateUserRequest | User data
    request: ...,
  } satisfies UsersPostRequest;

  try {
    const data = await api.usersPost(body);
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
| **request** | [WebCreateUserRequest](WebCreateUserRequest.md) | User data | |

### Return type

[**WebUserResponse**](WebUserResponse.md)

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

