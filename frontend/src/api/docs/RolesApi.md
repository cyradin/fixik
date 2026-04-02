# RolesApi

All URIs are relative to *http://localhost:8080/api*

| Method | HTTP request | Description |
|------------- | ------------- | -------------|
| [**rolesGet**](RolesApi.md#rolesget) | **GET** /roles | List roles |



## rolesGet

> WebListRolesResponse rolesGet()

List roles

Get all user roles

### Example

```ts
import {
  Configuration,
  RolesApi,
} from '';
import type { RolesGetRequest } from '';

async function example() {
  console.log("🚀 Testing  SDK...");
  const api = new RolesApi();

  try {
    const data = await api.rolesGet();
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

[**WebListRolesResponse**](WebListRolesResponse.md)

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

