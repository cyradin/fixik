
# WebCreateUserRequest


## Properties

Name | Type
------------ | -------------
`email` | string
`name` | string
`password` | string
`role` | string
`teamId` | number
`username` | string

## Example

```typescript
import type { WebCreateUserRequest } from ''

// TODO: Update the object below with actual values
const example = {
  "email": null,
  "name": null,
  "password": null,
  "role": null,
  "teamId": null,
  "username": null,
} satisfies WebCreateUserRequest

console.log(example)

// Convert the instance to a JSON string
const exampleJSON: string = JSON.stringify(example)
console.log(exampleJSON)

// Parse the JSON string back to an object
const exampleParsed = JSON.parse(exampleJSON) as WebCreateUserRequest
console.log(exampleParsed)
```

[[Back to top]](#) [[Back to API list]](../README.md#api-endpoints) [[Back to Model list]](../README.md#models) [[Back to README]](../README.md)


