# POPULATION PAGINATION API

An API Endpoint to access population data in Indonesia by district and sub-district. This project was created as a whole using the [Go](https://go.dev/) programming language with the [Fiber](https://docs.gofiber.io/) framework.

## Available features

- API Endpoint to request a response with a search and page parameter query which produces data for all districts in Indonesia which are provided in a paginated manner with a data limit per page of 10 data, there is data on the number of districts and the number of pages, as well as a list of all districts.
- API Endpoint to send a request which returns details of each district based on id. District details include a list of all sub-districts and the population of each sub-district.
