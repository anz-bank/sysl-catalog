

[Back](../README.md)


# BFF

## Integration Diagram
![](integration.svg)







## Application Index


| Application Name | Method | Source Location |
|----|----|----|
| DeliveryServer | [...](#DeliveryServer-...) | [https://github.com/anz-bank/sysl-catalog/blob/master/demo/simple2.sysl](https://github.com/anz-bank/sysl-catalog/blob/master/demo/simple2.sysl)|  
| IdentityServer | [Authenticate](#IdentityServer-Authenticate) | [https://github.com/anz-bank/sysl-catalog/blob/master/demo/simple2.sysl](https://github.com/anz-bank/sysl-catalog/blob/master/demo/simple2.sysl)|  
| IdentityServer | [CustomerProfile](#IdentityServer-CustomerProfile) | [https://github.com/anz-bank/sysl-catalog/blob/master/demo/simple2.sysl](https://github.com/anz-bank/sysl-catalog/blob/master/demo/simple2.sysl)|  
| IdentityServer | [NewCustomer](#IdentityServer-NewCustomer) | [https://github.com/anz-bank/sysl-catalog/blob/master/demo/simple2.sysl](https://github.com/anz-bank/sysl-catalog/blob/master/demo/simple2.sysl)|  
| IdentityServer | [UpdatePassword](#IdentityServer-UpdatePassword) | [https://github.com/anz-bank/sysl-catalog/blob/master/demo/simple2.sysl](https://github.com/anz-bank/sysl-catalog/blob/master/demo/simple2.sysl)|  
| OrderServer | [Order](#OrderServer-Order) | [https://github.com/anz-bank/sysl-catalog/blob/master/demo/simple2.sysl](https://github.com/anz-bank/sysl-catalog/blob/master/demo/simple2.sysl)|  
| OrderServer | [Review](#OrderServer-Review) | [https://github.com/anz-bank/sysl-catalog/blob/master/demo/simple2.sysl](https://github.com/anz-bank/sysl-catalog/blob/master/demo/simple2.sysl)|  
| OrderServer | [UpdateOrderStatus](#OrderServer-UpdateOrderStatus) | [https://github.com/anz-bank/sysl-catalog/blob/master/demo/simple2.sysl](https://github.com/anz-bank/sysl-catalog/blob/master/demo/simple2.sysl)|  
| PaymentServer | [Pay](#PaymentServer-Pay) | [https://github.com/anz-bank/sysl-catalog/blob/master/demo/simple2.sysl](https://github.com/anz-bank/sysl-catalog/blob/master/demo/simple2.sysl)|  
| ProductServer | [Menu](#ProductServer-Menu) | [https://github.com/anz-bank/sysl-catalog/blob/master/demo/simple2.sysl](https://github.com/anz-bank/sysl-catalog/blob/master/demo/simple2.sysl)|  




## Type Index


| Application Name | Type Name | Source Location |
|----|----|----|
| IdentityServer | [Customer](#IdentityServer.Customer) | [https://github.com/anz-bank/sysl-catalog/blob/master/demo/simple2.sysl](https://github.com/anz-bank/sysl-catalog/blob/master/demo/simple2.sysl)|
| IdentityServer | [NewCustomerRequest](#IdentityServer.NewCustomerRequest) | [https://github.com/anz-bank/sysl-catalog/blob/master/demo/simple2.sysl](https://github.com/anz-bank/sysl-catalog/blob/master/demo/simple2.sysl)|
| IdentityServer | [UnauthorizedError](#IdentityServer.UnauthorizedError) | [https://github.com/anz-bank/sysl-catalog/blob/master/demo/simple2.sysl](https://github.com/anz-bank/sysl-catalog/blob/master/demo/simple2.sysl)|
| OrderServer | [Order](#OrderServer.Order) | [https://github.com/anz-bank/sysl-catalog/blob/master/demo/simple2.sysl](https://github.com/anz-bank/sysl-catalog/blob/master/demo/simple2.sysl)|
| OrderServer | [OrderProduct](#OrderServer.OrderProduct) | [https://github.com/anz-bank/sysl-catalog/blob/master/demo/simple2.sysl](https://github.com/anz-bank/sysl-catalog/blob/master/demo/simple2.sysl)|
| OrderServer | [OrderRequest](#OrderServer.OrderRequest) | [https://github.com/anz-bank/sysl-catalog/blob/master/demo/simple2.sysl](https://github.com/anz-bank/sysl-catalog/blob/master/demo/simple2.sysl)|
| OrderServer | [OrderStatus](#OrderServer.OrderStatus) | [https://github.com/anz-bank/sysl-catalog/blob/master/demo/simple2.sysl](https://github.com/anz-bank/sysl-catalog/blob/master/demo/simple2.sysl)|
| PaymentServer | [PaymentType](#PaymentServer.PaymentType) | [https://github.com/anz-bank/sysl-catalog/blob/master/demo/simple2.sysl](https://github.com/anz-bank/sysl-catalog/blob/master/demo/simple2.sysl)|
| ProductServer | [Product](#ProductServer.Product) | [https://github.com/anz-bank/sysl-catalog/blob/master/demo/simple2.sysl](https://github.com/anz-bank/sysl-catalog/blob/master/demo/simple2.sysl)|
| ProductServer | [Products](#ProductServer.Products) | [https://github.com/anz-bank/sysl-catalog/blob/master/demo/simple2.sysl](https://github.com/anz-bank/sysl-catalog/blob/master/demo/simple2.sysl)|








# Applications





## Application DeliveryServer



- We are going to provide delivery service ASAP
 since our customers need it during COVID-19











### <a name=DeliveryServer-...></a>DeliveryServer ...


<details>
<summary>Sequence Diagram</summary>

![](DeliveryServer/....svg)
</details>

<details>
<summary>Request types</summary>


<span style="color:grey">No Request types</span>






</details>

<details>
<summary>Response types</summary>





<span style="color:grey">No Response Types</span>

</details>


---






## Application IdentityServer



- This server handles all the customer related endpoints
 including customer profile, password update, 
 customer authentication, etc.











### <a name=IdentityServer-Authenticate></a>IdentityServer Authenticate
this is a description of Authenticate

<details>
<summary>Sequence Diagram</summary>

![](IdentityServer/authenticate.svg)
</details>

<details>
<summary>Request types</summary>







![](primitive/stringemail.svg)



![](primitive/stringpassword.svg)



</details>

<details>
<summary>Response types</summary>





<span style="color:grey">No Response Types</span>

</details>


---





### <a name=IdentityServer-CustomerProfile></a>IdentityServer CustomerProfile


<details>
<summary>Sequence Diagram</summary>

![](IdentityServer/customerprofile.svg)
</details>

<details>
<summary>Request types</summary>







![](primitive/intcustomer_id.svg)



</details>

<details>
<summary>Response types</summary>






![](IdentityServer/customer.svg)




</details>


---





### <a name=IdentityServer-NewCustomer></a>IdentityServer NewCustomer


<details>
<summary>Sequence Diagram</summary>

![](IdentityServer/newcustomer.svg)
</details>

<details>
<summary>Request types</summary>







![](IdentityServer/newcustomerrequestreq.svg)



</details>

<details>
<summary>Response types</summary>






![](IdentityServer/customer.svg)




</details>


---





### <a name=IdentityServer-UpdatePassword></a>IdentityServer UpdatePassword


<details>
<summary>Sequence Diagram</summary>

![](IdentityServer/updatepassword.svg)
</details>

<details>
<summary>Request types</summary>







![](primitive/intcustomer_id.svg)



![](primitive/stringold.svg)



![](primitive/stringnew.svg)



</details>

<details>
<summary>Response types</summary>





<span style="color:grey">No Response Types</span>

</details>


---






## Application OrderServer



- This server handles all the order
 related endpoints.











### <a name=OrderServer-Order></a>OrderServer Order


<details>
<summary>Sequence Diagram</summary>

![](OrderServer/order.svg)
</details>

<details>
<summary>Request types</summary>







![](OrderServer/orderrequestreq.svg)



</details>

<details>
<summary>Response types</summary>





<span style="color:grey">No Response Types</span>

</details>


---





### <a name=OrderServer-Review></a>OrderServer Review


<details>
<summary>Sequence Diagram</summary>

![](OrderServer/review.svg)
</details>

<details>
<summary>Request types</summary>







![](primitive/intscore.svg)



![](primitive/stringcomment.svg)



</details>

<details>
<summary>Response types</summary>






![](OrderServer/order.svg)




</details>


---





### <a name=OrderServer-UpdateOrderStatus></a>OrderServer UpdateOrderStatus


<details>
<summary>Sequence Diagram</summary>

![](OrderServer/updateorderstatus.svg)
</details>

<details>
<summary>Request types</summary>







![](primitive/intorder_id.svg)



![](primitive/intstatus.svg)



</details>

<details>
<summary>Response types</summary>






![](OrderServer/order.svg)




</details>


---






## Application PaymentServer



- This server handles all the payment related endpoints.











### <a name=PaymentServer-Pay></a>PaymentServer Pay


<details>
<summary>Sequence Diagram</summary>

![](PaymentServer/pay.svg)
</details>

<details>
<summary>Request types</summary>


<span style="color:grey">No Request types</span>






</details>

<details>
<summary>Response types</summary>





<span style="color:grey">No Response Types</span>

</details>


---






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











<a name=IdentityServer.Customer></a><details>
<summary>IdentityServer.Customer</summary>

### IdentityServer.Customer

- Empty Type

![](IdentityServer/customersimple.svg)

[Full Diagram](IdentityServer/customer.svg)


#### Fields

| Field name | Type | Description |
|----|----|----|
| email | string | |
| first_name | string | |
| last_name | string | |
| middle_name | string | |
| phone | string | |


</details>
<a name=IdentityServer.NewCustomerRequest></a><details>
<summary>IdentityServer.NewCustomerRequest</summary>

### IdentityServer.NewCustomerRequest



![](IdentityServer/newcustomerrequestsimple.svg)

[Full Diagram](IdentityServer/newcustomerrequest.svg)


#### Fields

| Field name | Type | Description |
|----|----|----|
| email | string | |
| first_name | string | |
| last_name | string | |
| middle_name | string | |
| password | string | |
| phone | string | |


</details>
<a name=IdentityServer.UnauthorizedError></a><details>
<summary>IdentityServer.UnauthorizedError</summary>

### IdentityServer.UnauthorizedError



![](IdentityServer/unauthorizederrorsimple.svg)

[Full Diagram](IdentityServer/unauthorizederror.svg)


#### Fields

| Field name | Type | Description |
|----|----|----|
| error_msg | string | |


</details>




<a name=OrderServer.Order></a><details>
<summary>OrderServer.Order</summary>

### OrderServer.Order

- Customer order information

![](OrderServer/ordersimple.svg)

[Full Diagram](OrderServer/order.svg)


#### Fields

| Field name | Type | Description |
|----|----|----|
| id | int | |
| items | sequence of OrderProduct | |
| paid | bool | |
| review_comment | string | |
| review_score | int | |
| status | OrderStatus | |
| total_price | int | |


</details>
<a name=OrderServer.OrderProduct></a><details>
<summary>OrderServer.OrderProduct</summary>

### OrderServer.OrderProduct

- Order items

![](OrderServer/orderproductsimple.svg)

[Full Diagram](OrderServer/orderproduct.svg)


#### Fields

| Field name | Type | Description |
|----|----|----|
| comments | string | |
| product_id | int | |
| quantity | int | |


</details>
<a name=OrderServer.OrderRequest></a><details>
<summary>OrderServer.OrderRequest</summary>

### OrderServer.OrderRequest



![](OrderServer/orderrequestsimple.svg)

[Full Diagram](OrderServer/orderrequest.svg)


#### Fields

| Field name | Type | Description |
|----|----|----|
| order_id | int | |
| product_id | int | |
| quantity | int | |


</details>
<a name=OrderServer.OrderStatus></a><details>
<summary>OrderServer.OrderStatus</summary>

### OrderServer.OrderStatus



![](OrderServer/orderstatussimple.svg)

[Full Diagram](OrderServer/orderstatus.svg)



</details>




<a name=PaymentServer.PaymentType></a><details>
<summary>PaymentServer.PaymentType</summary>

### PaymentServer.PaymentType



![](PaymentServer/paymenttypesimple.svg)

[Full Diagram](PaymentServer/paymenttype.svg)



</details>




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

