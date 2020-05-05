
[Back](../README.md)
# Package ServerPackage

## Service Index
| Service Name | Method |
----|----
Server | [Authenticate](#Server-Authenticate) |
Server | [GET/testRestPathParamPrimitive/{primitiveID}](#Server-GET/testRestPathParamPrimitive/{primitiveID}) |
Server | [GET/testRestQueryParam](#Server-GET/testRestQueryParam) |
Server | [GET/testRestQueryParamPrimitive](#Server-GET/testRestQueryParamPrimitive) |
Server | [GET/testRestURLParam/{id}](#Server-GET/testRestURLParam/{id}) |

## Database Index
| Database Name |
----|
| [RelModel](#Database-RelModel) |

[Types](#Types)

## Integration diagram

![](ServerPackage_integration.svg)

---



---




## Server

- this is a comment for Server





## Server Authenticate

this is a description of Authenticate

### Sequence Diagram
![](ServerAuthenticate.svg)

### Request types


![](ServerAuthenticatedata-model-parameter0.svg)






### Response types


![](ServerAuthenticatedata-model-response0.svg)






## Server GET /testRestPathParamPrimitive/{primitiveID}



### Sequence Diagram
![](ServerGETtestRestPathParamPrimitive{primitiveID}.svg)

### Request types






![](ServerGETtestRestPathParamPrimitive{primitiveID}data-model-path-parameter0.svg)


### Response types


![](ServerGETtestRestPathParamPrimitive{primitiveID}data-model-response0.svg)






## Server GET /testRestQueryParam



### Sequence Diagram
![](ServerGETtestRestQueryParam.svg)

### Request types




![](ServerGETtestRestQueryParamdata-model-query-parameter0.svg)


![](ServerGETtestRestQueryParamdata-model-query-parameter1.svg)




### Response types


![](ServerGETtestRestQueryParamdata-model-response0.svg)






## Server GET /testRestQueryParamPrimitive



### Sequence Diagram
![](ServerGETtestRestQueryParamPrimitive.svg)

### Request types




![](ServerGETtestRestQueryParamPrimitivedata-model-query-parameter0.svg)




### Response types


![](ServerGETtestRestQueryParamPrimitivedata-model-response0.svg)






## Server GET /testRestURLParam/{id}



### Sequence Diagram
![](ServerGETtestRestURLParam{id}.svg)

### Request types






![](ServerGETtestRestURLParam{id}data-model-path-parameter0.svg)


### Response types


![](ServerGETtestRestURLParam{id}data-model-response0.svg)


---



## Database RelModel

![](RelModeldb.svg)


## Types
<table>
<tr>
<th>App Name</th>
<th>Diagram</th>
<th>Comment</th>
<th>Full Diagram</th>
</tr>


<tr>
<td>

MegaDatabase.<br>Empty
</td>
<td>

![](SimpleEmptydata-model1.svg)
</td>
<td> 

 
</td>
<td>

[Link](Full-Emptydata-model1.svg)
</td>
</tr>
<tr>
<td>

MegaDatabase.<br>Money
</td>
<td>

![](SimpleMoneydata-model1.svg)
</td>
<td> 

 
</td>
<td>

[Link](Full-Moneydata-model1.svg)
</td>
</tr>
<tr>
<td>

Server.<br>Request
</td>
<td>

![](SimpleRequestdata-model1.svg)
</td>
<td> 

 
</td>
<td>

[Link](Full-Requestdata-model1.svg)
</td>
</tr>
<tr>
<td>

Server.<br>Response
</td>
<td>

![](SimpleResponsedata-model1.svg)
</td>
<td> 

 
</td>
<td>

[Link](Full-Responsedata-model1.svg)
</td>
</tr>
</table>
