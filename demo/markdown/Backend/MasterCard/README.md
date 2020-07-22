

[Back](../README.md)


# MasterCard

## Integration Diagram
![](integration.svg)







## Application Index


| Application Name | Method | Source Location |
|----|----|----|
| MasterCard | [POST /pay](#MasterCard-POSTpay) | [https://github.com/anz-bank/sysl-catalog/blob/master/demo/simple.yaml](https://github.com/anz-bank/sysl-catalog/blob/master/demo/simple.yaml)|  




## Type Index


| Application Name | Type Name | Source Location |
|----|----|----|
| MasterCard | [AustralianState](#MasterCard.AustralianState) | [https://github.com/anz-bank/sysl-catalog/blob/master/demo/simple.yaml](https://github.com/anz-bank/sysl-catalog/blob/master/demo/simple.yaml)|
| MasterCard | [SimpleObj](#MasterCard.SimpleObj) | [https://github.com/anz-bank/sysl-catalog/blob/master/demo/simple.yaml](https://github.com/anz-bank/sysl-catalog/blob/master/demo/simple.yaml)|
| MasterCard | [SimpleObj2](#MasterCard.SimpleObj2) | [https://github.com/anz-bank/sysl-catalog/blob/master/demo/simple.yaml](https://github.com/anz-bank/sysl-catalog/blob/master/demo/simple.yaml)|








# Applications





## Application MasterCard



- No description.











### <a name=MasterCard-POSTpay></a>MasterCard POST /pay


<details>
<summary>Sequence Diagram</summary>

![](MasterCard/postpay.svg)
</details>

<details>
<summary>Request types</summary>



<span style="color:grey">No Request types</span>







</details>

<details>
<summary>Response types</summary>






![](MasterCard/simpleobj.svg)




</details>


---





# Types







<a name=MasterCard.AustralianState></a><details>
<summary>MasterCard.AustralianState</summary>

### MasterCard.AustralianState



![](MasterCard/australianstatesimple.svg)

[Full Diagram](MasterCard/australianstate.svg)



</details>
<a name=MasterCard.SimpleObj></a><details>
<summary>MasterCard.SimpleObj</summary>

### MasterCard.SimpleObj



![](MasterCard/simpleobjsimple.svg)

[Full Diagram](MasterCard/simpleobj.svg)


#### Fields

| Field name | Type | Description |
|----|----|----|
| name | string | |


</details>
<a name=MasterCard.SimpleObj2></a><details>
<summary>MasterCard.SimpleObj2</summary>

### MasterCard.SimpleObj2



![](MasterCard/simpleobj2simple.svg)

[Full Diagram](MasterCard/simpleobj2.svg)


#### Fields

| Field name | Type | Description |
|----|----|----|
| name | SimpleObj | |


</details>


<div class="footer">

