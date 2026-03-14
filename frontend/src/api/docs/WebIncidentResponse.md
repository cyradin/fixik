
# WebIncidentResponse


## Properties

Name | Type
------------ | -------------
`description` | string
`id` | number
`impact` | [WebDictEntityShort](WebDictEntityShort.md)
`priority` | [WebDictEntityShort](WebDictEntityShort.md)
`status` | [WebDictEntityShort](WebDictEntityShort.md)
`title` | string

## Example

```typescript
import type { WebIncidentResponse } from ''

// TODO: Update the object below with actual values
const example = {
  "description": null,
  "id": null,
  "impact": null,
  "priority": null,
  "status": null,
  "title": null,
} satisfies WebIncidentResponse

console.log(example)

// Convert the instance to a JSON string
const exampleJSON: string = JSON.stringify(example)
console.log(exampleJSON)

// Parse the JSON string back to an object
const exampleParsed = JSON.parse(exampleJSON) as WebIncidentResponse
console.log(exampleParsed)
```

[[Back to top]](#) [[Back to API list]](../README.md#api-endpoints) [[Back to Model list]](../README.md#models) [[Back to README]](../README.md)


