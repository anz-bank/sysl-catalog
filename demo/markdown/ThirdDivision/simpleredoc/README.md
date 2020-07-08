

[Back](../README.md) | [Chat with us](https://anzoss.slack.com/messages/sysl-catalog/) | [New bug or feature request](https://github.com/anz-bank/sysl-catalog/issues/new)


# simpleredoc

## Integration Diagram
![](integration.svg)







## Application Index


| Application Name | Method | Source Location |
|----|----|----|
| simpleredoc | [GET /test](#simpleredoc-GETtest) | [https://github.com/anz-bank/sysl-catalog/blob/master/demo/simple.yaml](https://github.com/anz-bank/sysl-catalog/blob/master/demo/simple.yaml)|  




## Type Index


| Application Name | Type Name | Source Location |
|----|----|----|
| simpleredoc | [AustralianState](#simpleredoc.AustralianState) | [https://github.com/anz-bank/sysl-catalog/blob/master/demo/simple.yaml](https://github.com/anz-bank/sysl-catalog/blob/master/demo/simple.yaml)|
| simpleredoc | [SimpleObj](#simpleredoc.SimpleObj) | [https://github.com/anz-bank/sysl-catalog/blob/master/demo/simple.yaml](https://github.com/anz-bank/sysl-catalog/blob/master/demo/simple.yaml)|
| simpleredoc | [SimpleObj2](#simpleredoc.SimpleObj2) | [https://github.com/anz-bank/sysl-catalog/blob/master/demo/simple.yaml](https://github.com/anz-bank/sysl-catalog/blob/master/demo/simple.yaml)|








# Applications





## Application simpleredoc



- No description.











### <a name=simpleredoc-GETtest></a>simpleredoc GET /test


<details>
<summary>Sequence Diagram</summary>

![](simpleredoc/gettest.svg)
</details>

<details>
<summary>Request types</summary>



<span style="color:grey">No Request types</span>







</details>

<details>
<summary>Response types</summary>






![](simpleredoc/simpleobj.svg)




</details>


---





# Types







<a name=simpleredoc.AustralianState></a><details>
<summary>simpleredoc.AustralianState</summary>

### simpleredoc.AustralianState



![](simpleredoc/australianstatesimple.svg)

[Full Diagram](simpleredoc/australianstate.svg)



</details>
<a name=simpleredoc.SimpleObj></a><details>
<summary>simpleredoc.SimpleObj</summary>

### simpleredoc.SimpleObj



![](simpleredoc/simpleobjsimple.svg)

[Full Diagram](simpleredoc/simpleobj.svg)


#### Fields

| Field name | Type | Description |
|----|----|----|
| name | string | |


</details>
<a name=simpleredoc.SimpleObj2></a><details>
<summary>simpleredoc.SimpleObj2</summary>

### simpleredoc.SimpleObj2



![](simpleredoc/simpleobj2simple.svg)

[Full Diagram](simpleredoc/simpleobj2.svg)


#### Fields

| Field name | Type | Description |
|----|----|----|
| name | SimpleObj | |


</details>


<div class="footer">

