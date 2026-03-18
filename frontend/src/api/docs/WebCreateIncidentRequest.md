
# WebCreateIncidentRequest


## Properties

Name | Type
------------ | -------------
`authorId` | number
`description` | string
`priorityId` | number
`statusId` | number
`teamId` | number
`title` | string
`userId` | number

## Example

```typescript
import type { WebCreateIncidentRequest } from ''

// TODO: Update the object below with actual values
const example = {
  "authorId": null,
  "description": null,
  "priorityId": null,
  "statusId": null,
  "teamId": null,
  "title": null,
  "userId": null,
} satisfies WebCreateIncidentRequest

console.log(example)

// Convert the instance to a JSON string
const exampleJSON: string = JSON.stringify(example)
console.log(exampleJSON)

// Parse the JSON string back to an object
const exampleParsed = JSON.parse(exampleJSON) as WebCreateIncidentRequest
console.log(exampleParsed)
```

[[Back to top]](#) [[Back to API list]](../README.md#api-endpoints) [[Back to Model list]](../README.md#models) [[Back to README]](../README.md)


