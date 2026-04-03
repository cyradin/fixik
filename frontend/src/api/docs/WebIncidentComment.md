
# WebIncidentComment


## Properties

Name | Type
------------ | -------------
`author` | [WebUserResponse](WebUserResponse.md)
`createdAt` | string
`id` | number
`incidentId` | number
`text` | string
`updatedAt` | string

## Example

```typescript
import type { WebIncidentComment } from ''

// TODO: Update the object below with actual values
const example = {
  "author": null,
  "createdAt": null,
  "id": null,
  "incidentId": null,
  "text": null,
  "updatedAt": null,
} satisfies WebIncidentComment

console.log(example)

// Convert the instance to a JSON string
const exampleJSON: string = JSON.stringify(example)
console.log(exampleJSON)

// Parse the JSON string back to an object
const exampleParsed = JSON.parse(exampleJSON) as WebIncidentComment
console.log(exampleParsed)
```

[[Back to top]](#) [[Back to API list]](../README.md#api-endpoints) [[Back to Model list]](../README.md#models) [[Back to README]](../README.md)


