Tax Calculator Example
======================
This is a sample API to generate tax foreach user.
Generate using [https://github.com/syroegkin/swagger-markdown](https://github.com/syroegkin/swagger-markdown).

**Version:** 1.0


**License:** [Apache 2.0](http://www.apache.org/licenses/LICENSE-2.0.html)

### /login
---
##### ***POST***
**Summary:** Login account

**Description:** Login using username and password

**Parameters**

| Name | Located in | Description | Required | Schema |
| ---- | ---------- | ----------- | -------- | ---- |
| user | body | user info | Yes | [reqpayload.Login](#reqpayload.login) |

**Responses**

| Code | Description | Schema |
| ---- | ----------- | ------ |
| 200 | OK | [respayload.Login](#respayload.login) |
| 400 | Bad Request | [respayload.Error](#respayload.error) |
| 401 | Unauthorized | [respayload.Error](#respayload.error) |
| 404 | Not Found | [respayload.Error](#respayload.error) |
| 422 | Unprocessable Entity | [respayload.Error](#respayload.error) |

### /register
---
##### ***POST***
**Summary:** Register new account

**Description:** Register new account

**Parameters**

| Name | Located in | Description | Required | Schema |
| ---- | ---------- | ----------- | -------- | ---- |
| user | body | user info | Yes | [reqpayload.Register](#reqpayload.register) |

**Responses**

| Code | Description | Schema |
| ---- | ----------- | ------ |
| 200 | OK | [respayload.Register](#respayload.register) |
| 400 | Bad Request | [respayload.Error](#respayload.error) |
| 422 | Unprocessable Entity | [respayload.Error](#respayload.error) |

### /tax
---
##### ***GET***
**Summary:** Get taxes related to current user

**Description:** Get taxes related to current user. Tax calculation is based on following calculation rule: 1. Food and Beverage: 10% of Price, for example if the price is 1000 then the tax is 100, hence the amount is 1100. 2. Tobacco: 10 + (2% of Price), for example if the price is 1000 then the tax is 10 + (2% * 1000) = 10 + 20 = 30, hence the amount is 1030. 3. Entertainment: if the price is equal or more than 100 is 1% of (Price - 100), otherwise is free. For instance, if the price is 150, then the tax is 1% * (150-100) = 1% * 50 = 0.5, hence the final amount is 150.5.

**Parameters**

| Name | Located in | Description | Required | Schema |
| ---- | ---------- | ----------- | -------- | ---- |
| Authentication-Token | header | Authentication-Token your-token | Yes | string |

**Responses**

| Code | Description | Schema |
| ---- | ----------- | ------ |
| 200 | OK | [respayload.TaxesForCurrentUser](#respayload.taxesforcurrentuser) |
| 400 | Bad Request | [respayload.Error](#respayload.error) |
| 422 | Unprocessable Entity | [respayload.Error](#respayload.error) |

##### ***POST***
**Summary:** Add tax record to your account

**Description:** Add tax record to your account

**Parameters**

| Name | Located in | Description | Required | Schema |
| ---- | ---------- | ----------- | -------- | ---- |
| Authentication-Token | header | Authentication-Token your-token | Yes | string |
| tax | body | tax info | Yes | [reqpayload.CreateNewTax](#reqpayload.createnewtax) |

**Responses**

| Code | Description | Schema |
| ---- | ----------- | ------ |
| 200 | OK | [respayload.Tax](#respayload.tax) |
| 400 | Bad Request | [respayload.Error](#respayload.error) |
| 422 | Unprocessable Entity | [respayload.Error](#respayload.error) |

### Models
---

### reqpayload.CreateNewTax  

| Name | Type | Description | Required |
| ---- | ---- | ----------- | -------- |
| name | string |  | Yes |
| price | integer |  | Yes |
| tax_code | integer |  | Yes |

### reqpayload.Login  

| Name | Type | Description | Required |
| ---- | ---- | ----------- | -------- |
| password | string |  | Yes |
| username | string |  | Yes |

### reqpayload.Register  

| Name | Type | Description | Required |
| ---- | ---- | ----------- | -------- |
| password | string |  | Yes |
| username | string |  | Yes |

### respayload.Error  

| Name | Type | Description | Required |
| ---- | ---- | ----------- | -------- |
| error_code | string |  | No |
| http_status_code | integer |  | No |
| message | string |  | No |

### respayload.Login  

| Name | Type | Description | Required |
| ---- | ---- | ----------- | -------- |
| authentication_token | string |  | No |
| user | [respayload.User](#respayload.user) |  | No |

### respayload.Register  

| Name | Type | Description | Required |
| ---- | ---- | ----------- | -------- |
| authentication_token | string |  | No |
| user | [respayload.User](#respayload.user) |  | No |

### respayload.Tax  

| Name | Type | Description | Required |
| ---- | ---- | ----------- | -------- |
| amount | string |  | No |
| name | string |  | No |
| price | integer |  | No |
| refundable | boolean |  | No |
| tax | string |  | No |
| tax_code | integer |  | No |
| type | string |  | No |

### respayload.TaxesForCurrentUser  

| Name | Type | Description | Required |
| ---- | ---- | ----------- | -------- |
| grand_total | string |  | No |
| price_sub_total | integer |  | No |
| tax_sub_total | string |  | No |
| taxes | [ [respayload.Tax](#respayload.tax) ] |  | No |

### respayload.User  

| Name | Type | Description | Required |
| ---- | ---- | ----------- | -------- |
| id | integer |  | No |
| username | string |  | No |