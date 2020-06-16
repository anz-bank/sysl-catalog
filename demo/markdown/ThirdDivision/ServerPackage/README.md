

[Back](../README.md)


# ServerPackage

## Integration Diagram
![](demo/markdown/integration.svg)








## Database Index
| Database Application Name  | Source Location |
----|----
[RelModel](#Database-RelModel) | [../../../../demo/simple2.sysl](../../../../demo/simple2.sysl)|  


## Application Index
| Application Name | Method | Source Location |
----|----|----
Server | [Authenticate](#Server-Authenticate) | [../../../../demo/simple2.sysl](../../../../demo/simple2.sysl)|  
Server | [GET /testRestPathParamPrimitive/{primitiveID}](#Server-GETtestRestPathParamPrimitive{primitiveID}) | [../../../../demo/simple2.sysl](../../../../demo/simple2.sysl)|  
Server | [GET /testRestQueryParam](#Server-GETtestRestQueryParam) | [../../../../demo/simple2.sysl](../../../../demo/simple2.sysl)|  
Server | [GET /testRestQueryParamPrimitive](#Server-GETtestRestQueryParamPrimitive) | [../../../../demo/simple2.sysl](../../../../demo/simple2.sysl)|  
Server | [GET /testRestURLParam/{id}](#Server-GETtestRestURLParam{id}) | [../../../../demo/simple2.sysl](../../../../demo/simple2.sysl)|  
Server | [GET /testReturnNil](#Server-GETtestReturnNil) | [../../../../demo/simple2.sysl](../../../../demo/simple2.sysl)|  

## Type Index
| Application Name | Type Name | Source Location |
----|----|----
Server | [Empty](#Server.Empty) | [../../../../demo/simple2.sysl](../../../../demo/simple2.sysl)|
Server | [Request](#Server.Request) | [../../../../demo/simple2.sysl](../../../../demo/simple2.sysl)|
Server | [Response](#Server.Response) | [../../../../demo/simple2.sysl](../../../../demo/simple2.sysl)|



# Databases



<details>
<summary>Database RelModel</summary>


![](demo/markdown/RelModel/types.svg)
</details>




# Applications







## Application Server

- this is a comment for Server







### Server Authenticate
this is a description of Authenticate

<details>
<summary>Sequence Diagram</summary>

![](demo/markdown/Server/authenticate.svg)
</details>

<details>
<summary>Request types</summary>

#### Request types





![](demo/markdown/Server/request.svg)



</details>
<details>
<summary>Response types</summary>

#### Response types





![](demo/markdown/Server/response.svg)



</details>

---





### Server GETtestRestPathParamPrimitive{primitiveID}


<details>
<summary>Sequence Diagram</summary>

![](demo/markdown/Server/gettestrestpathparamprimitive{primitiveid}.svg)
</details>

<details>
<summary>Request types</summary>

#### Request types








#### Path Parameter

![](demo/markdown/primitive/stringsimple.svg)



</details>
<details>
<summary>Response types</summary>

#### Response types





![](demo/markdown/Server/response.svg)



</details>

---





### Server GETtestRestQueryParam


<details>
<summary>Sequence Diagram</summary>

![](demo/markdown/Server/gettestrestqueryparam.svg)
</details>

<details>
<summary>Request types</summary>

#### Request types










#### Query Parameter

![](demo/markdown/Server/request.svg)



#### Query Parameter

![](demo/markdown/Server/request.svg)

</details>
<details>
<summary>Response types</summary>

#### Response types





![](demo/markdown/Server/response.svg)



</details>

---





### Server GETtestRestQueryParamPrimitive


<details>
<summary>Sequence Diagram</summary>

![](demo/markdown/Server/gettestrestqueryparamprimitive.svg)
</details>

<details>
<summary>Request types</summary>

#### Request types










#### Query Parameter

![](demo/markdown/primitive/stringsimple.svg)

</details>
<details>
<summary>Response types</summary>

#### Response types





![](demo/markdown/Server/response.svg)



</details>

---





### Server GETtestRestURLParam{id}


<details>
<summary>Sequence Diagram</summary>

![](demo/markdown/Server/gettestresturlparam{id}.svg)
</details>

<details>
<summary>Request types</summary>

#### Request types








#### Path Parameter

![](demo/markdown/Server/request.svg)



</details>
<details>
<summary>Response types</summary>

#### Response types





![](demo/markdown/Server/response.svg)



</details>

---





### Server GETtestReturnNil


<details>
<summary>Sequence Diagram</summary>

![](demo/markdown/Server/gettestreturnnil.svg)
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

![](demo/markdown/Server/emptysimple.svg)

[Full Diagram](demo/markdown/Server/empty.svg)

#### Fields

| Field name | Type | Description |
|----|----|----|

</details>
<details>
<summary>Server.Request</summary>

### Server.Request

- Request Request Request

![](demo/markdown/Server/requestsimple.svg)

[Full Diagram](demo/markdown/Server/request.svg)

#### Fields

| Field name | Type | Description |
|----|----|----|
| query | sequence of Response | |

</details>
<details>
<summary>Server.Response</summary>

### Server.Response

- Response Response Response

![](demo/markdown/Server/responsesimple.svg)

[Full Diagram](demo/markdown/Server/response.svg)

#### Fields

| Field name | Type | Description |
|----|----|----|
| balance | MegaDatabase.Empty | |
| query | MegaDatabase.Money | |

</details>

<div class="footer">

