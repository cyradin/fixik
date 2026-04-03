
# WebIncidentResponse


## Properties

Name | Type
------------ | -------------
`author` | [WebUserResponse](WebUserResponse.md)
`commentsCount` | number
`createdAt` | string
`description` | string
`id` | number
`priority` | [WebDictEntityShort](WebDictEntityShort.md)
`status` | [WebDictEntityShort](WebDictEntityShort.md)
`team` | [WebDictEntityShort](WebDictEntityShort.md)
`title` | string
`updatedAt` | string
`user` | [WebUserResponse](WebUserResponse.md)

## Example

```typescript
import type { WebIncidentResponse } from ''

// TODO: Update the object below with actual values
const example = {
  "author": null,
  "commentsCount": null,
  "createdAt": null,
  "description": null,
  "id": null,
  "priority": null,
  "status": null,
  "team": null,
  "title": null,
  "updatedAt": null,
  "user": null,
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


