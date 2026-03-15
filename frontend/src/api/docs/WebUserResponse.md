
# WebUserResponse


## Properties

Name | Type
------------ | -------------
`email` | string
`id` | number
`name` | string
`role` | string
`teamId` | number
`username` | string

## Example

```typescript
import type { WebUserResponse } from ''

// TODO: Update the object below with actual values
const example = {
  "email": null,
  "id": null,
  "name": null,
  "role": null,
  "teamId": null,
  "username": null,
} satisfies WebUserResponse

console.log(example)

// Convert the instance to a JSON string
const exampleJSON: string = JSON.stringify(example)
console.log(exampleJSON)

// Parse the JSON string back to an object
const exampleParsed = JSON.parse(exampleJSON) as WebUserResponse
console.log(exampleParsed)
```

[[Back to top]](#) [[Back to API list]](../README.md#api-endpoints) [[Back to Model list]](../README.md#models) [[Back to README]](../README.md)


