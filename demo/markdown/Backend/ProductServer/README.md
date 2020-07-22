

[Back](../README.md)


# ProductServer

## Integration Diagram
![](integration.svg)







## Application Index


| Application Name | Method | Source Location |
|----|----|----|
| ProductServer | [Menu](#ProductServer-Menu) | [https://github.com/anz-bank/sysl-catalog/blob/master/demo/sizzle.sysl](https://github.com/anz-bank/sysl-catalog/blob/master/demo/sizzle.sysl)|  




## Type Index


| Application Name | Type Name | Source Location |
|----|----|----|
| ProductServer | [Product](#ProductServer.Product) | [https://github.com/anz-bank/sysl-catalog/blob/master/demo/sizzle.sysl](https://github.com/anz-bank/sysl-catalog/blob/master/demo/sizzle.sysl)|
| ProductServer | [Products](#ProductServer.Products) | [https://github.com/anz-bank/sysl-catalog/blob/master/demo/sizzle.sysl](https://github.com/anz-bank/sysl-catalog/blob/master/demo/sizzle.sysl)|








# Applications





## Application ProductServer



- This server handles all the product
 related endpoints.











### <a name=ProductServer-Menu></a>ProductServer Menu


<details>
<summary>Sequence Diagram</summary>

![](ProductServer/menu.svg)
</details>

<details>
<summary>Request types</summary>


<span style="color:grey">No Request types</span>






</details>

<details>
<summary>Response types</summary>






![](ProductServer/products.svg)




</details>


---





# Types







<a name=ProductServer.Product></a><details>
<summary>ProductServer.Product</summary>

### ProductServer.Product

- Product information

![](ProductServer/productsimple.svg)

[Full Diagram](ProductServer/product.svg)


#### Fields

| Field name | Type | Description |
|----|----|----|
| details | string | |
| id | int | |
| image | string | |
| name | string | |
| price | int | |


</details>
<a name=ProductServer.Products></a><details>
<summary>ProductServer.Products</summary>

### ProductServer.Products



![](ProductServer/productssimple.svg)

[Full Diagram](ProductServer/products.svg)



</details>


<div class="footer">

