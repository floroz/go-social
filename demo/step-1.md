## Problem Statement

We have a working OpenAPI spec for a Social Network app and we want to introduce a new resource: “Likes”

### Initialize memory bank

- Start by prompting to initialize the bank

>Prompt (ACT):
>
>initialize memory bank

- Empty templates are scaffolded

### Provide the necessary context on the project and the requirements

- Stay in PLAN mode and begin providing the necessary context

>Prompt (PLAN):
>
>This project contains a @/README.md that provides an overview over what it is. There is a backend in Golang and a Frontend in React.
>
>One important aspect of this project is that it follows OpenAPI specification to design and provide standards around how the endpoints are built and what contracts are exposed to the user, which in our case for now is only the Frontend.
>
>However, in the future we can consider exposing these endpoints to clients too (we are already following a versioning convention to avoid breaking changes).
>
>The specs are defined in @/openapi - we write partials that are then recomposed together to generate a final spec, the latest one is @/openapi/openapi-bundled.yaml
>
>We have not yet implemented a strategy on how to produce and persist different bundles that can represent snapshots of our API spec at different points in time. This will be needed and is a missing feature to implement in the future.
>
>A more immediate requirement we have is the definition of a new endpoint. This endpoint should expose "likes".
>
>I have shared a lot of context with you so far, some of which I believe should be saved in our memory bank.
>Regarding the new endpoint implementation, let's begin a discussion on how it should be structured - keeping in mind we can lean on the existing patterns and conventions of the existing openapi spec.

- The model will generate a plan - review it and if happy - select ACT mode.
- The model will begin updating the memory bank.
- Review the files that have been created.

### Review initial plan draft in the memory bank

- We start in [systemPatterns.md](../memory-bank/systemPatterns.md) where we see some critical findings discovered by Cline, particularly in "Key Technical Decisions"
  - Cline discovered that we have codegen for API types
  - Cline discovered that we use the /v1/ path versioning
  - Cline discovered we use structured response for Success and Errors

We didn't directly provide them with this information, but it's now available context for the model, and extremely relevant to our task.

- In [techContext.md](../memory-bank/techContext.md) Cline identified the usage of some important scripts from the Makefile.

- In [activeContext.md](../memory-bank/activeContext.md) what Cline considers to be the plan:

- We see some redundancy across [projectbrief.md](../memory-bank/projectbrief.md) and [productContext.md](../memory-bank/productContext.md) - we will address that if we start seeing results that conflict with the content of those files.

```
## Next Steps

- Complete updates for `systemPatterns.md`, `techContext.md`, and `progress.md`.
- Discuss the requirements and design for the new "likes" endpoint.
- Define the OpenAPI specification for the "likes" endpoint.
- Implement the backend logic for "likes".
- Implement the frontend UI and logic for "likes".
```

- ***Scope creep***: The model is assuming we want to implement the full functionality - including UI and BE logic in our plan.

- The model needs guidance to avoid tackling backend (BE) and frontend (FE) development before finalizing the API specification, similar to how a junior developer might try to do everything simultaneously. 


### Let's Chunk the Plan

> Prompt (PLAN):
> The context you saved and identified appears mostly correct. However we need to adjust the Next steps:
>
> - Complete updates for systemPatterns.md, techContext.md, and progress.md.
> - Discuss the requirements and design for the new "likes" endpoint.
> - Define the OpenAPI specification for the "likes" endpoint.
> - Implement the backend logic for "likes".
> - Implement the frontend UI and logic for "likes".
>
> The implementation of the feature should be out of scope - we are only interested in the endpoint API Spec at this stage, and we consider this work "Done" once we have a valid OpenAPI spec that we can load to our current Swagger UI. The Swagger UI is already functional and will automatically pick updates to the spec.
>
> Implement versioning of the bundled spec is also out of scope for now.
>
> Let's update our plan to reflect these requirements.

We are now ready for Step 2.