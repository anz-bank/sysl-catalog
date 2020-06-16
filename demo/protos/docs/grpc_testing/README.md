

[Back](../README.md)


# grpc_testing

## Integration Diagram
![](../../docs/images/grpc_testing-integration.svg)







## Application Index
| Application Name | Method | Source Location |
----|----|----
Bar | [AnotherEndpoint](#Bar-AnotherEndpoint) | [../../simple.proto](../../simple.proto)|  
Foo | [thisEndpoint](#Foo-thisEndpoint) | [../../simple.proto](../../simple.proto)|  

## Type Index
| Application Name | Type Name | Source Location |
----|----|----
grpc_testing | [Money](#grpc_testing.Money) | [../../simple.proto](../../simple.proto)|
grpc_testing | [Request](#grpc_testing.Request) | [../../simple.proto](../../simple.proto)|
grpc_testing | [Response](#grpc_testing.Response) | [../../simple.proto](../../simple.proto)|




# Applications





## Application Bar

- This is a comment before Bar







### Bar AnotherEndpoint
this is a comment before Bar.AnotherEndpoint

<details>
<summary>Sequence Diagram</summary>

![](../../docs/images/grpc_testing-Bar-anotherendpoint.svg)
</details>

<details>
<summary>Request types</summary>

#### Request types





![](../../docs/images/grpc_testing-grpc_testing-request.svg)



</details>
<details>
<summary>Response types</summary>

#### Response types






![](../../docs/images/grpc_testing-grpc_testing-response.svg)



</details>

---






## Application Foo

- This is a comment before Foo







### Foo thisEndpoint


<details>
<summary>Sequence Diagram</summary>

![](../../docs/images/grpc_testing-Foo-thisendpoint.svg)
</details>

<details>
<summary>Request types</summary>

#### Request types





![](../../docs/images/grpc_testing-grpc_testing-request.svg)



</details>
<details>
<summary>Response types</summary>

#### Response types






![](../../docs/images/grpc_testing-grpc_testing-response.svg)



</details>

---







# Types








<details>
<summary>grpc_testing.Money</summary>

### grpc_testing.Money

- 

![](../../docs/images/grpc_testing-grpc_testing-moneysimple.svg)

[Full Diagram](../../docs/images/grpc_testing-grpc_testing-money.svg)

#### Fields

| Field name | Type | Description |
|----|----|----|
| nanos | int | |
| units | int | |

</details>
<details>
<summary>grpc_testing.Request</summary>

### grpc_testing.Request

- 

![](../../docs/images/grpc_testing-grpc_testing-requestsimple.svg)

[Full Diagram](../../docs/images/grpc_testing-grpc_testing-request.svg)

#### Fields

| Field name | Type | Description |
|----|----|----|
| query | string | |

</details>
<details>
<summary>grpc_testing.Response</summary>

### grpc_testing.Response

- 

![](../../docs/images/grpc_testing-grpc_testing-responsesimple.svg)

[Full Diagram](../../docs/images/grpc_testing-grpc_testing-response.svg)

#### Fields

| Field name | Type | Description |
|----|----|----|
| query | string | |

</details>

<div class="footer">

