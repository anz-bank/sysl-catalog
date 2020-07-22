

[Back](../README.md)


# IdentityServer

## Integration Diagram
![](integration.svg)







## Application Index


| Application Name | Method | Source Location |
|----|----|----|
| IdentityServer | [Authenticate](#IdentityServer-Authenticate) | [https://github.com/anz-bank/sysl-catalog/blob/master/demo/sizzle.sysl](https://github.com/anz-bank/sysl-catalog/blob/master/demo/sizzle.sysl)|  
| IdentityServer | [CustomerProfile](#IdentityServer-CustomerProfile) | [https://github.com/anz-bank/sysl-catalog/blob/master/demo/sizzle.sysl](https://github.com/anz-bank/sysl-catalog/blob/master/demo/sizzle.sysl)|  
| IdentityServer | [NewCustomer](#IdentityServer-NewCustomer) | [https://github.com/anz-bank/sysl-catalog/blob/master/demo/sizzle.sysl](https://github.com/anz-bank/sysl-catalog/blob/master/demo/sizzle.sysl)|  
| IdentityServer | [UpdatePassword](#IdentityServer-UpdatePassword) | [https://github.com/anz-bank/sysl-catalog/blob/master/demo/sizzle.sysl](https://github.com/anz-bank/sysl-catalog/blob/master/demo/sizzle.sysl)|  




## Type Index


| Application Name | Type Name | Source Location |
|----|----|----|
| IdentityServer | [Customer](#IdentityServer.Customer) | [https://github.com/anz-bank/sysl-catalog/blob/master/demo/sizzle.sysl](https://github.com/anz-bank/sysl-catalog/blob/master/demo/sizzle.sysl)|
| IdentityServer | [NewCustomerRequest](#IdentityServer.NewCustomerRequest) | [https://github.com/anz-bank/sysl-catalog/blob/master/demo/sizzle.sysl](https://github.com/anz-bank/sysl-catalog/blob/master/demo/sizzle.sysl)|
| IdentityServer | [UnauthorizedError](#IdentityServer.UnauthorizedError) | [https://github.com/anz-bank/sysl-catalog/blob/master/demo/sizzle.sysl](https://github.com/anz-bank/sysl-catalog/blob/master/demo/sizzle.sysl)|








# Applications





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





# Types







<a name=IdentityServer.Customer></a><details>
<summary>IdentityServer.Customer</summary>

### IdentityServer.Customer

- This contains all information relating to a customer

![](IdentityServer/customersimple.svg)

[Full Diagram](IdentityServer/customer.svg)


#### Fields

| Field name | Type | Description |
|----|----|----|
| email | string | |
| first_name | string | |
| last_name | string | |
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


<div class="footer">

