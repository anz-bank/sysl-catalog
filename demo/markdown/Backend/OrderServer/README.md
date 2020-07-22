

[Back](../README.md)


# OrderServer

## Integration Diagram
![](integration.svg)







## Application Index


| Application Name | Method | Source Location |
|----|----|----|
| OrderServer | [Order](#OrderServer-Order) | [https://github.com/anz-bank/sysl-catalog/blob/master/demo/sizzle.sysl](https://github.com/anz-bank/sysl-catalog/blob/master/demo/sizzle.sysl)|  
| OrderServer | [Review](#OrderServer-Review) | [https://github.com/anz-bank/sysl-catalog/blob/master/demo/sizzle.sysl](https://github.com/anz-bank/sysl-catalog/blob/master/demo/sizzle.sysl)|  
| OrderServer | [UpdateOrderStatus](#OrderServer-UpdateOrderStatus) | [https://github.com/anz-bank/sysl-catalog/blob/master/demo/sizzle.sysl](https://github.com/anz-bank/sysl-catalog/blob/master/demo/sizzle.sysl)|  




## Type Index


| Application Name | Type Name | Source Location |
|----|----|----|
| OrderServer | [Order](#OrderServer.Order) | [https://github.com/anz-bank/sysl-catalog/blob/master/demo/sizzle.sysl](https://github.com/anz-bank/sysl-catalog/blob/master/demo/sizzle.sysl)|
| OrderServer | [OrderProduct](#OrderServer.OrderProduct) | [https://github.com/anz-bank/sysl-catalog/blob/master/demo/sizzle.sysl](https://github.com/anz-bank/sysl-catalog/blob/master/demo/sizzle.sysl)|
| OrderServer | [OrderRequest](#OrderServer.OrderRequest) | [https://github.com/anz-bank/sysl-catalog/blob/master/demo/sizzle.sysl](https://github.com/anz-bank/sysl-catalog/blob/master/demo/sizzle.sysl)|
| OrderServer | [OrderStatus](#OrderServer.OrderStatus) | [https://github.com/anz-bank/sysl-catalog/blob/master/demo/sizzle.sysl](https://github.com/anz-bank/sysl-catalog/blob/master/demo/sizzle.sysl)|








# Applications





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





# Types







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


<div class="footer">

