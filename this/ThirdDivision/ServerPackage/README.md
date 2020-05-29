

[Back](../README.md)


# ServerPackage

## Integration Diagram
![](integration.svg)








## Database Index
| Database Application Name  | Source Location |
----|----
[RelModel](#Database-RelModel) | [../../../demo/simple2.sysl](../../../demo/simple2.sysl)|  


## Application Index
| Application Name | Method | Source Location |
----|----|----
Server | [Authenticate](#Server-Authenticate) | [../../../demo/simple2.sysl](../../../demo/simple2.sysl)|  
Server | [GET /testRestPathParamPrimitive/{primitiveID}](#Server-GETtestRestPathParamPrimitive{primitiveID}) | [../../../demo/simple2.sysl](../../../demo/simple2.sysl)|  
Server | [GET /testRestQueryParam](#Server-GETtestRestQueryParam) | [../../../demo/simple2.sysl](../../../demo/simple2.sysl)|  
Server | [GET /testRestQueryParamPrimitive](#Server-GETtestRestQueryParamPrimitive) | [../../../demo/simple2.sysl](../../../demo/simple2.sysl)|  
Server | [GET /testRestURLParam/{id}](#Server-GETtestRestURLParam{id}) | [../../../demo/simple2.sysl](../../../demo/simple2.sysl)|  
Server | [GET /testReturnNil](#Server-GETtestReturnNil) | [../../../demo/simple2.sysl](../../../demo/simple2.sysl)|  

## Type Index
| Application Name | Type Name | Source Location |
----|----|----
Server | [Empty](#Server.Empty) | [../../../demo/simple2.sysl](../../../demo/simple2.sysl)|
Server | [Request](#Server.Request) | [../../../demo/simple2.sysl](../../../demo/simple2.sysl)|
Server | [Response](#Server.Response) | [../../../demo/simple2.sysl](../../../demo/simple2.sysl)|



# Databases



<details>
<summary>Database RelModel</summary>


![](RelModel/types.svg)
</details>




# Applications







## Application Server

- this is a comment for Server





### Server Authenticate
this is a description of Authenticate

<details>
<summary>Sequence Diagram</summary>

![](Server/authenticate.svg)
</details>

<details>
<summary>Request types</summary>

#### Request types





![](Server/request.svg)



</details>
<details>
<summary>Response types</summary>

#### Response types





![](Server/response.svg)



</details>

---





### Server GETtestRestPathParamPrimitive{primitiveID}


<details>
<summary>Sequence Diagram</summary>

![](Server/gettestrestpathparamprimitive{primitiveid}.svg)
</details>

<details>
<summary>Request types</summary>

#### Request types








#### Path Parameter

![](primitive/stringsimple.svg)



</details>
<details>
<summary>Response types</summary>

#### Response types





![](Server/response.svg)



</details>

---





### Server GETtestRestQueryParam


<details>
<summary>Sequence Diagram</summary>

![](Server/gettestrestqueryparam.svg)
</details>

<details>
<summary>Request types</summary>

#### Request types










#### Query Parameter

![](Server/request.svg)



#### Query Parameter

![](Server/request.svg)

</details>
<details>
<summary>Response types</summary>

#### Response types





![](Server/response.svg)



</details>

---





### Server GETtestRestQueryParamPrimitive


<details>
<summary>Sequence Diagram</summary>

![](Server/gettestrestqueryparamprimitive.svg)
</details>

<details>
<summary>Request types</summary>

#### Request types










#### Query Parameter

![](primitive/stringsimple.svg)

</details>
<details>
<summary>Response types</summary>

#### Response types





![](Server/response.svg)



</details>

---





### Server GETtestRestURLParam{id}


<details>
<summary>Sequence Diagram</summary>

![](Server/gettestresturlparam{id}.svg)
</details>

<details>
<summary>Request types</summary>

#### Request types








#### Path Parameter

![](Server/request.svg)



</details>
<details>
<summary>Response types</summary>

#### Response types





![](Server/response.svg)



</details>

---





### Server GETtestReturnNil


<details>
<summary>Sequence Diagram</summary>

![](Server/gettestreturnnil.svg)
</details>

<details>
<summary>Request types</summary>

#### Request types







</details>
<details>
<summary>Response types</summary>

#### Response types



No Response Types


</details>

---




# Types





<details>
<summary>Server.Empty</summary>

### Server.Empty

- Empty Empty Empty

![](Server/emptysimple.svg)

[Full Diagram](Server/empty.svg)

#### Fields

| Field name | Type | Description |
|----|----|----|

</details>
<details>
<summary>Server.Request</summary>

### Server.Request

- Request Request Request

![](Server/requestsimple.svg)

[Full Diagram](Server/request.svg)

#### Fields

| Field name | Type | Description |
|----|----|----|
| query | sequence of Response | |

</details>
<details>
<summary>Server.Response</summary>

### Server.Response

- Response Response Response

![](Server/responsesimple.svg)

[Full Diagram](Server/response.svg)

#### Fields

| Field name | Type | Description |
|----|----|----|
| balance | MegaDatabase.Empty | |
| query | MegaDatabase.Money | |

</details>

<div class="footer">

