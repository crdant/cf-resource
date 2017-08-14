# Cloud Foundry Route Resource

An output only resource (at the moment) that will create/map/unmap routes for cloud foundry.

## Source Configuration

* `api`: *Required.* The address of the Cloud Controller in the Cloud Foundry
  deployment.
* `username`: *Required.* The username used to authenticate.
* `password`: *Required.* The password used to authenticate.
* `organization`: *Required.* The organization to push the application to.
* `space`: *Required.* The space to push the application to.
* `skip_cert_check`: *Optional.* Check the validity of the CF SSL cert.
  Defaults to `false`.

## Behaviour

### `out`: Create, Map, or Unmap Routes

Pushes an application to the Cloud Foundry detailed in the source
configuration. A [manifest][cf-manifests] that describes the application must
be specified.

[cf-manifests]: http://docs.cloudfoundry.org/devguide/deploy-apps/manifest.html

Application string   `json:"application"`
Create      []string `json:"create"`
RandomPort  bool     `json:"randomPort"`
Map         []string `json:"unmap"`
Unmap       []string `json:"map"`

#### Parameters

*N.B. one of `map`, `create`, or `unmap` must be provided*

* `create`: *Required.* One or more routes to create, in the same format as a  route specified in a [Cloud Foundry manifest](https://docs.cloudfoundry.org/devguide/deploy-apps/manifest.html).
* `map`: *Optional.* One or more routes to map to an application, in the same format as a  route specified in a [Cloud Foundry manifest](https://docs.cloudfoundry.org/devguide/deploy-apps/manifest.html). Requires that `application` be specified.
* `unmap`: *Optional.* One or more routes to unmap from an application, in the same format as a  route specified in a [Cloud Foundry manifest](https://docs.cloudfoundry.org/devguide/deploy-apps/manifest.html). Requires that `application` be specified.
* `application`: *Optional.* The name of an application for which to map/unmap the route. Required if `map` or `unmap` is provided.
* `randomPort`: *Optional.* Use a random port when creating a TCP route. Specify only the domain for `create` when using this option.
