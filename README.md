# MyRamen 서비스 REST API 명세서

## 기본 정보
- Base URL: `https://myramen-api.injun.dev/api/v1`
- 응답 형식: JSON
- 문자 인코딩: UTF-8

## 환경 변수
- `ADMIN_PASSWORD`: 관리자 작업을 위한 비밀번호

## API 엔드포인트

### 1. 라면 구매 요청
> 참고: 라면 3개 이상 구매 시 자동으로 200원 할인 쿠폰 발급

**요청 정보:**
- URL: `/orders`
- 메소드: `POST`
- Content-Type: `application/json`

**요청 본문 (Request Body):**
```json
{
  "name": "홍길동",                // 주문자 이름 (필수)
  "accountNumber": "123-456-789",  // 계좌번호 (필수)
  "quantity": 5,                   // 라면 수량 (필수, 양의 정수)
  "spicyLevel": 3,                 // 매운맛 정도 (1: 순한맛 ~ 5: 매운맛, 기본값: 3)
  "deliveryOption": "PICKUP_4F",   // 배달 방식 (PICKUP_4F: 4층 자습실 픽업, PICKUP_LAUNDRY: 세탁실 픽업, DELIVERY: 배달)
  "options": {                     // 추가 옵션
    "chopsticks": true,            // 젓가락 포함 여부
    "hotWaterDelivery": false,     // 뜨거운 물 배달 서비스 (+500원)
    "cookingService": false        // 조리 서비스 (+500원)
  },
  "couponId": "c78910"             // 사용할 쿠폰 ID (선택 사항)
}
```

**응답:**
- 상태 코드: `201 Created` (성공 시)

**응답 본문 (Response Body):**
```json
{
  "orderId": "o12345",            // 주문 고유 ID
  "name": "홍길동",               // 주문자 이름
  "accountNumber": "123-456-789", // 계좌번호
  "quantity": 5,                  // 라면 수량
  "spicyLevel": 3,                // 매운맛 정도
  "deliveryOption": "PICKUP_4F",  // 배달 방식
  "options": {                    // 추가 옵션
    "chopsticks": true,
    "hotWaterDelivery": false,
    "cookingService": false
  },
  "totalPrice": 19800,            // 총 가격 (라면 가격 + 추가 옵션 가격 - 쿠폰 할인)
  "appliedCoupon": {              // 적용된 쿠폰 정보 (쿠폰 사용 시에만 포함)
    "couponId": "c78910",         // 쿠폰 고유 ID
    "discount": 200               // 할인 금액
  },
  "newCoupon": {                  // 발급된 새 쿠폰 정보 (quantity가 3 이상일 때만 포함)
    "couponId": "c78912",         // 쿠폰 고유 ID
    "discount": 200,              // 할인 금액
    "expiryDate": "2025-06-08T23:59:59Z" // 만료일 (발급일로부터 30일)
  },
  "status": "PENDING"             // 주문 상태 (기본값: PENDING)
}
```

**오류 응답:**
- 상태 코드: `400 Bad Request` (요청 데이터 오류)

```json
{
  "error": "INVALID_REQUEST",
  "message": "구매 요청 정보가 유효하지 않습니다."
}
```

- 상태 코드: `400 Bad Request` (유효하지 않은 쿠폰)

```json
{
  "error": "INVALID_COUPON",
  "message": "사용할 수 없는 쿠폰입니다. (이미 사용되었거나 만료됨)"
}
```

### 2. 라면 주문 조회(상태 확인)

**요청 정보:**
- URL: `/orders/{orderId}`
- 메소드: `GET`

**경로 파라미터:**
- `orderId`: 주문 고유 ID

**응답:**
- 상태 코드: `200 OK` (성공 시)

**응답 본문 (Response Body):**
```json
{
  "orderId": "o12345",
  "name": "홍길동",
  "accountNumber": "123-456-789",
  "quantity": 5,
  "spicyLevel": 3,
  "deliveryOption": "PICKUP_4F",
  "options": {
    "chopsticks": true,
    "hotWaterDelivery": false,
    "cookingService": false
  },
  "totalPrice": 19800,
  "appliedCoupon": {
    "couponId": "c78910",
    "discount": 200
  },
  "status": "PENDING"     // 주문 상태 (PENDING, PAID, COOKING, READY, DELIVERING, DELIVERED)
}
```

**오류 응답:**
- 상태 코드: `404 Not Found` (주문 ID를 찾을 수 없는 경우)

```json
{
  "error": "NOT_FOUND",
  "message": "해당 주문을 찾을 수 없습니다."
}
```

### 3. 모든 주문 목록 조회(관리자용)

**요청 정보:**
- URL: `/admin/orders`
- 메소드: `GET`
- Headers:
  - `X-Admin-Password`: 환경 변수로 설정된 관리자 비밀번호

**응답:**
- 상태 코드: `200 OK`

**응답 본문 (Response Body):**
```json
{
  "orders": [
    {
      "orderId": "o12345",
      "name": "홍길동",
      "accountNumber": "123-456-789",
      "quantity": 5,
      "spicyLevel": 3,
      "deliveryOption": "PICKUP_4F",
      "options": {
        "chopsticks": true,
        "hotWaterDelivery": false,
        "cookingService": false
      },
      "totalPrice": 19800,
      "appliedCoupon": {
        "couponId": "c78910",
        "discount": 200
      },
      "status": "PENDING"
    },
    {
      "orderId": "o12346",
      "name": "김철수",
      "accountNumber": "987-654-321",
      "quantity": 3,
      "spicyLevel": 5,
      "deliveryOption": "DELIVERY",
      "options": {
        "chopsticks": true,
        "hotWaterDelivery": true,
        "cookingService": true
      },
      "totalPrice": 13000,
      "status": "PAID"
    }
    // ... 더 많은 주문 데이터
  ]
}
```

**오류 응답:**
- 상태 코드: `401 Unauthorized` (비밀번호가 일치하지 않는 경우)

```json
{
  "error": "UNAUTHORIZED",
  "message": "관리자 인증에 실패했습니다."
}
```

### 4. 특정 주문 상태 변경(관리자용)

**요청 정보:**
- URL: `/admin/orders/{orderId}/status`
- 메소드: `PUT`
- Content-Type: `application/json`
- Headers:
  - `X-Admin-Password`: 환경 변수로 설정된 관리자 비밀번호

**경로 파라미터:**
- `orderId`: 주문 고유 ID

**요청 본문 (Request Body):**
```json
{
  "status": "COOKING"    // 변경할 상태 (PENDING, PAID, COOKING, READY, DELIVERING, DELIVERED)
}
```

**응답:**
- 상태 코드: `200 OK`

**응답 본문 (Response Body):**
```json
{
  "orderId": "o12345",
  "name": "홍길동",
  "accountNumber": "123-456-789",
  "quantity": 5,
  "spicyLevel": 3,
  "deliveryOption": "PICKUP_4F",
  "options": {
    "chopsticks": true,
    "hotWaterDelivery": false,
    "cookingService": false
  },
  "totalPrice": 19800,
  "appliedCoupon": {
    "couponId": "c78910",
    "discount": 200
  },
  "status": "COOKING"    // 변경된 상태
}
```

**오류 응답:**
- 상태 코드: `401 Unauthorized` (비밀번호가 일치하지 않는 경우)

```json
{
  "error": "UNAUTHORIZED",
  "message": "관리자 인증에 실패했습니다."
}
```

- 상태 코드: `404 Not Found` (주문 ID를 찾을 수 없는 경우)

```json
{
  "error": "NOT_FOUND",
  "message": "해당 주문을 찾을 수 없습니다."
}
```

- 상태 코드: `400 Bad Request` (유효하지 않은 상태)

```json
{
  "error": "INVALID_STATUS",
  "message": "유효하지 않은 주문 상태입니다."
}
```

### 5. 쿠폰 조회

**요청 정보:**
- URL: `/coupons/{couponId}`
- 메소드: `GET`

**경로 파라미터:**
- `couponId`: 쿠폰 고유 ID

**응답:**
- 상태 코드: `200 OK` (성공 시)

**응답 본문 (Response Body):**
```json
{
  "couponId": "c78910",
  "discount": 200,
  "expiryDate": "2025-06-08T23:59:59Z",
  "isUsed": false,
  "issuedAt": "2025-05-08T14:30:00Z"
}
```

**오류 응답:**
- 상태 코드: `404 Not Found` (쿠폰 ID를 찾을 수 없는 경우)

```json
{
  "error": "NOT_FOUND",
  "message": "해당 쿠폰을 찾을 수 없습니다."
}
```

### 6. 모든 유효한 쿠폰 목록 조회 (관리자용)

**요청 정보:**
- URL: `/admin/coupons`
- 메소드: `GET`
- Headers:
  - `X-Admin-Password`: 환경 변수로 설정된 관리자 비밀번호

**응답:**
- 상태 코드: `200 OK`

**응답 본문 (Response Body):**
```json
{
  "coupons": [
    {
      "couponId": "c78910",
      "discount": 200,
      "expiryDate": "2025-06-08T23:59:59Z",
      "isUsed": false,
      "issuedAt": "2025-05-08T14:30:00Z"
    },
    {
      "couponId": "c78911",
      "discount": 200,
      "expiryDate": "2025-06-10T23:59:59Z",
      "isUsed": false,
      "issuedAt": "2025-05-10T10:15:00Z"
    }
  ]
}
```

**오류 응답:**
- 상태 코드: `401 Unauthorized` (비밀번호가 일치하지 않는 경우)

```json
{
  "error": "UNAUTHORIZED",
  "message": "관리자 인증에 실패했습니다."
}
```

## 데이터 모델

### 주문(Order)
| 필드 | 타입 | 설명 |
|------|------|------|
| orderId | String | 주문 고유 ID |
| name | String | 주문자 이름 |
| accountNumber | String | 계좌번호 |
| quantity | Integer | 라면 수량 |
| spicyLevel | Integer | 매운맛 정도 (1: 순한맛 ~ 5: 매운맛) |
| deliveryOption | String | 배달 방식 (PICKUP_4F, PICKUP_LAUNDRY, DELIVERY) |
| options | Object | 추가 옵션 (chopsticks, hotWaterDelivery, cookingService) |
| totalPrice | Integer | 총 가격 |
| status | String | 주문 상태 |
| appliedCoupon | Object | 적용된 쿠폰 정보 (쿠폰 사용 시에만 포함) |
| newCoupon | Object | 새로 발급된 쿠폰 정보 (quantity가 3 이상일 때만 포함) |

### 쿠폰(Coupon)
| 필드 | 타입 | 설명 |
|------|------|------|
| couponId | String | 쿠폰 고유 ID |
| discount | Integer | 할인 금액 (200원) |
| expiryDate | DateTime | 만료일 (발급일로부터 30일) |
| isUsed | Boolean | 사용 여부 |
| issuedAt | DateTime | 발급일 |

### 주문 상태
| 상태 | 설명 |
|------|------|
| PENDING | 주문 대기 중 (결제 전) |
| PAID | 결제 완료 |
| COOKING | 조리 중 |
| READY | 픽업 준비 완료 |
| DELIVERING | 배달 중 |
| DELIVERED | 배달/픽업 완료 |

### 가격 정보
| 항목 | 가격 |
|------|------|
| 신라면 기본 가격 | 4,000원 |
| 뜨거운 물 배달 서비스 | +500원 |
| 조리 서비스 | +500원 |

### 프로모션 정보
| 프로모션 | 설명 |
|---------|------|
| 3개 이상 구매 쿠폰 | 라면 3개 이상 구매 시 다음 주문에 사용 가능한 200원 할인 쿠폰 자동 발급 (유효기간: 30일) |

## 오류 코드
| 코드 | HTTP 상태 | 설명 |
|------|-----------|------|
| INVALID_REQUEST | 400 | 요청 데이터가 유효하지 않음 |
| INVALID_STATUS | 400 | 유효하지 않은 주문 상태 |
| INVALID_COUPON | 400 | 유효하지 않은 쿠폰 (이미 사용됨/만료됨) |
| UNAUTHORIZED | 401 | 관리자 인증 실패 |
| NOT_FOUND | 404 | 리소스를 찾을 수 없음 |
