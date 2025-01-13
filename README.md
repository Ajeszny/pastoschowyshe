# pastoschowyshe

## Endpoints:
- get_pasta_list
  
  GET
  args: page, pagesize, *if none are provided fetches full list*

  [{id, name, tags}]
  
- get_pasta/{id}
  
  GET
  fetches a pasta

- add_pasta

  POST
  adds a pasta to a list
  Requires authorization in form of Authorization: "Bearer {token}"

- login

  POST
  logs in with password, returns token
