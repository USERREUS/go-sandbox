## Query

### Find all orders
```
{
  findAllOrders {
    order_code
    date
    order_item {
      product_code
      name
      count
      cost
    }
  }
}
```

### Find order by code
```
{
  findOrder(code: "<code>") {
    order_code
    date
    order_item {
      product_code
      name
      count
      cost
    }
  }
}

```

### Find all products
```
{
  findAllProducts {
    code
    name
    weight
    description
  }
}
```

### Find product by code
```
{
  findProduct(code: "<code>") {
    code
    name
    weight
    description
  }
}
```


## Mutation

### Create order with slice order items
```
mutation {
  addOrder (
    orderInput: {
      items: [
        {
          code: "<code>", 
          name: "<name>", 
          count: <count>, 
          cost:  <cost>
        },...
      ]}) 
  {
    order_code
    date
    order_item {
      product_code
      name
      count
      cost
    }
  }
}
```