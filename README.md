# azdevops

`azdevops` is a *low level* SDK for Azure DevOps, based on the `azcore`/`azidentity` libraries.

## Goals

- Consistent and powerful client implementation based on `azcore`
- Multiple credential choices, based on `azidentity`
- Supports both ADO and ADO server (aka. TFS)

## API behaviors

Useful reference: https://learn.microsoft.com/en-us/azure/devops/extend/develop/work-with-urls?toc=%2Fazure%2Fdevops%2Fmarketplace-extensibility%2Ftoc.json&view=azure-devops&tabs=http

### `_apis/ResourceAreas`

Users can call `GET` to this API path, either behind an organization name(e.g. `/my-org/_apis/ResourceAreas`), or as the base path (`/_apis/ResourceAreas`). This endpoint supports anonymous call, no credential is needed. The organization scoped endpoint can return more resource areas based on the organization configuration.

The result of calling, e.g. organization level, resource areas are always the same, no matter which API version is used (e.g. no version, 5.0-preview.1, 7.2-preview.1).

### `_apis`

Users can call `OPTIONS` to this API path. The domain can be the primary domain (i.e. `dev.azure.com` ) or any area specific domains (e.g. `feeds.dev.azure.com`).

The result will be different based on different domain you are calling to.

Not all domains returned from `_apis/ResourceAreas` can successfully serve the `_apis` request:

- Some of them will just fail:
    - 400: Our services aren't available right now
    - 404: The resource cannot be found
- Some of them will return 200, but no content is returned

This might indicate these endpoints are not for API consumption, but might be a web UI or something else. We can safely ignore them.

For those domains that can successful serve the `_apis` request, the result are always the same, no matter which API version is used.
