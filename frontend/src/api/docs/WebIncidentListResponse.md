
# WebIncidentListResponse


## Properties

Name | Type
------------ | -------------
`items` | [Array&lt;WebIncidentResponse&gt;](WebIncidentResponse.md)
`pagination` | [WebPagination](WebPagination.md)

## Example

```typescript
import type { WebIncidentListResponse } from ''

// TODO: Update the object below with actual values
const example = {
  "items": null,
  "pagination": null,
} satisfies WebIncidentListResponse

console.log(example)

// Convert the instance to a JSON string
const exampleJSON: string = JSON.stringify(example)
console.log(exampleJSON)

// Parse the JSON string back to an object
const exampleParsed = JSON.parse(exampleJSON) as WebIncidentListResponse
console.log(exampleParsed)
```

[[Back to top]](#) [[Back to API list]](../README.md#api-endpoints) [[Back to Model list]](../README.md#models) [[Back to README]](../README.md)


