# Changes

## [1.5.0](https://www.github.com/googleapis/google-cloud-go/compare/v1.4.0...v1.5.0) (2021-02-24)


### Features

* **firestore:** add opencensus tracing support  ([#2942](https://www.github.com/googleapis/google-cloud-go/issues/2942)) ([257f322](https://www.github.com/googleapis/google-cloud-go/commit/257f322e68b75765bd316ccefed5461d4df538a0))


### Bug Fixes

* **firestore:** address a missing branch in watch.stop() error remapping ([#3643](https://www.github.com/googleapis/google-cloud-go/issues/3643)) ([89ad55d](https://www.github.com/googleapis/google-cloud-go/commit/89ad55d72f79995a68f9c2ed1cd9b5ba50009d6d))

## [1.4.0](https://www.github.com/googleapis/google-cloud-go/compare/firestore/v1.3.0...v1.4.0) (2020-12-03)


### Features

* **firestore:** support "!=" and "not-in" query operators ([#3207](https://www.github.com/googleapis/google-cloud-go/issues/3207)) ([5c44019](https://www.github.com/googleapis/google-cloud-go/commit/5c440192105fe3e9b5dd1b584118b309113935e3)), closes [/firebase.google.com/support/release-notes/js#version_7210_-_september_17_2020](https://www.github.com/googleapis//firebase.google.com/support/release-notes/js/issues/version_7210_-_september_17_2020)

## v1.3.0

- Add support for LimitToLast feature for queries. This allows
  a query to return the final N results. See docs
  [here](https://firebase.google.com/docs/reference/js/firebase.database.Query#limittolast).
- Add support for FieldTransformMinimum and FieldTransformMaximum.
- Add exported SetGoogleClientInfo method.
- Various updates to autogenerated clients.

## v1.2.0

- Deprecate v1beta1 client.
- Fix serverTimestamp docs.
- Add missing operators to query docs.
- Make document IDs 20 alpha-numeric characters. Previously, there could be more
  than 20 non-alphanumeric characters, which broke some users. See
  https://github.com/googleapis/google-cloud-go/issues/1715.
- Various updates to autogenerated clients.

## v1.1.1

- Fix bug in CollectionGroup query validation.

## v1.1.0

- Add support for `in` and `array-contains-any` query operators.

## v1.0.0

This is the first tag to carve out firestore as its own module. See:
https://github.com/golang/go/wiki/Modules#is-it-possible-to-add-a-module-to-a-multi-module-repository.