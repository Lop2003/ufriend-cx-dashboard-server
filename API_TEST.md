# 🚀 UFriend CX Dashboard - API Testing Guide

เอกสารนี้เป็นคู่มือสำหรับทดสอบ API ทั้งหมดของ **UFriend CX Dashboard Server** (พัฒนาด้วย Go Fiber, MongoDB) โดยมีทั้งคำอธิบาย, Query Parameters, โครงสร้าง Request Body รวมถึงคำสั่ง `curl` และรูปแบบ **REST Client (VS Code)** เพื่อให้สามารถทดสอบได้ทันที

---

## 📌 1. ข้อมูลการตั้งค่าเริ่มต้น (Configuration & Setup)

* **Base URL:** `http://localhost:3000` (หรือพอร์ตอื่นที่ตั้งค่าไว้ใน Env `PORT`)
* **Default Database:** MongoDB (`ufriend_cx`)

---

## 🩺 2. ระบบตรวจสอบสถานะ (Health Check)

ตรวจสอบความพร้อมใช้งานของระบบ Server

### **Get System Health**
* **Method:** `GET`
* **Path:** `/health`
* **คำสั่ง curl:**
  ```bash
  curl -X GET http://localhost:3000/health
  ```
* **ตัวอย่าง Success Response (200 OK):**
  ```json
  {
    "status": "ok"
  }
  ```

---

## 👥 3. ข้อมูลลูกค้า (Customer API)

ใช้จัดการดึงข้อมูลลูกค้าและสถิติภาพรวม

### **3.1 ดึงรายการลูกค้าทั้งหมด (List Customers)**
ดึงรายชื่อลูกค้าทั้งหมด สามารถใช้ Query Parameters ในการกรองข้อมูลได้
* **Method:** `GET`
* **Path:** `/api/customers`
* **Query Parameters:**
  * `branch` (string, optional) - กรองตามสาขา เช่น `Bangkok`, `Nonthaburi`
  * `status` (string, optional) - กรองตามสถานะ เช่น `active`, `completed`, `overdue`
* **คำสั่ง curl:**
  ```bash
  # ดึงลูกค้าทั้งหมด
  curl -X GET "http://localhost:3000/api/customers"

  # กรองเฉพาะลูกค้าสาขา Bangkok ที่มีสถานะ active
  curl -X GET "http://localhost:3000/api/customers?branch=Bangkok&status=active"
  ```
* **ตัวอย่าง Success Response (200 OK):**
  ```json
  {
    "success": true,
    "message": "customers retrieved",
    "data": [
      {
        "id": "60d5ecb862b53b3b44747209",
        "name": "สมชาย ใจดี",
        "phone": "0812345678",
        "product": "UFriend Special Plan",
        "branch": "Bangkok",
        "plan_months": 12,
        "status": "active",
        "created_at": "2026-05-19T05:21:54Z"
      }
    ]
  }
  ```

---

### **3.2 ดึงรายละเอียดลูกค้ารายบุคคล (Get Customer Detail)**
ดึงข้อมูลลูกค้ารายคน พร้อมกับรายการ Feedback และ Follow-up ทั้งหมดของลูกค้ารายนั้น
* **Method:** `GET`
* **Path:** `/api/customers/:id` (แทนที่ `:id` ด้วย MongoDB ObjectID ของลูกค้า)
* **คำสั่ง curl:**
  ```bash
  curl -X GET http://localhost:3000/api/customers/60d5ecb862b53b3b44747209
  ```
* **ตัวอย่าง Success Response (200 OK):**
  ```json
  {
    "success": true,
    "message": "customer retrieved",
    "data": {
      "id": "60d5ecb862b53b3b44747209",
      "name": "สมชาย ใจดี",
      "phone": "0812345678",
      "product": "UFriend Special Plan",
      "branch": "Bangkok",
      "plan_months": 12,
      "status": "active",
      "created_at": "2026-05-19T05:21:54Z",
      "feedbacks": [
        {
          "id": "60d5ecd062b53b3b4474720a",
          "rating": 5,
          "comment": "บริการดีเยี่ยมมากครับ พนักงานพูดจาไพเราะ",
          "category": "service",
          "sentiment": "positive",
          "created_at": "2026-05-19T06:12:00Z"
        }
      ],
      "follow_ups": [
        {
          "id": "60d5ece262b53b3b4474720b",
          "type": "feedback_reply",
          "note": "โทรไปขอบคุณลูกค้าและเสนอส่วนลดรอบถัดไป",
          "status": "done",
          "created_at": "2026-05-19T07:30:00Z"
        }
      ]
    }
  }
  ```

---

### **3.3 สถิติภาพรวม (Get Summary Statistics)**
ดึงข้อมูลสถิติภาพรวมของระบบสำหรับหน้า Dashboard (ยอดรวมลูกค้า, เรตติ้งเฉลี่ย, จำนวนสถานะต่างๆ)
* **Method:** `GET`
* **Path:** `/api/stats/summary`
* **คำสั่ง curl:**
  ```bash
  curl -X GET http://localhost:3000/api/stats/summary
  ```
* **ตัวอย่าง Success Response (200 OK):**
  ```json
  {
    "success": true,
    "message": "summary retrieved",
    "data": {
      "total_customers": 150,
      "avg_rating": 4.25,
      "overdue_count": 12,
      "active_count": 108,
      "completed_count": 30
    }
  }
  ```

---

### **3.4 สถิติแยกตามสาขา (Get Stats by Branch)**
ดึงข้อมูลสรุปสถิติจำนวนลูกค้า เรตติ้งเฉลี่ย และยอดค้างติดตามแยกตามสาขา
* **Method:** `GET`
* **Path:** `/api/stats/by-branch`
* **คำสั่ง curl:**
  ```bash
  curl -X GET http://localhost:3000/api/stats/by-branch
  ```
* **ตัวอย่าง Success Response (200 OK):**
  ```json
  {
    "success": true,
    "message": "branch stats retrieved",
    "data": [
      {
        "branch": "Bangkok",
        "customer_count": 75,
        "avg_rating": 4.5,
        "overdue_count": 4
      },
      {
        "branch": "Nonthaburi",
        "customer_count": 75,
        "avg_rating": 4.0,
        "overdue_count": 8
      }
    ]
  }
  ```

---

## 💬 4. ข้อมูลฟีดแบ็ก (Feedback API)

ใช้บันทึกและดูฟีดแบ็กจากลูกค้า

### **4.1 ดึงรายการ Feedback ทั้งหมด (List Feedbacks)**
* **Method:** `GET`
* **Path:** `/api/feedbacks`
* **Query Parameters:**
  * `category` (string, optional) - กรองตามหมวดหมู่: `service`, `payment`, `product`, `branch`
  * `rating` (int, optional) - กรองตามคะแนน (1-5)
* **คำสั่ง curl:**
  ```bash
  # ดึง Feedback ทั้งหมด
  curl -X GET "http://localhost:3000/api/feedbacks"

  # กรองเฉพาะหมวดหมู่บริการ (service) ที่ได้คะแนน 5 ดาว
  curl -X GET "http://localhost:3000/api/feedbacks?category=service&rating=5"
  ```
* **ตัวอย่าง Success Response (200 OK):**
  ```json
  {
    "success": true,
    "message": "feedbacks retrieved",
    "data": [
      {
        "id": "60d5ecd062b53b3b4474720a",
        "customer_id": "60d5ecb862b53b3b44747209",
        "rating": 5,
        "comment": "บริการดีเยี่ยมมากครับ พนักงานพูดจาไพเราะ",
        "category": "service",
        "sentiment": "positive",
        "created_at": "2026-05-19T06:12:00Z"
      }
    ]
  }
  ```

---

### **4.2 สร้าง Feedback ใหม่ (Create Feedback)**
* **Method:** `POST`
* **Path:** `/api/feedbacks`
* **Headers:** `Content-Type: application/json`
* **เงื่อนไขการ Validation (กฎความปลอดภัย):**
  * `customer_id` (ต้องกรอก และต้องเป็นรูปแบบ ObjectID ที่ถูกต้อง)
  * `rating` (ต้องเป็นตัวเลขจำนวนเต็มระหว่าง `1` ถึง `5` เท่านั้น)
  * `comment` (ต้องกรอกข้อความ)
  * `category` (ต้องกรอก และต้องเป็นคำใดคำหนึ่งดังต่อไปนี้: `service`, `payment`, `product`, `branch`)
* **โครงสร้าง Request Body (JSON):**
  ```json
  {
    "customer_id": "60d5ecb862b53b3b44747209",
    "rating": 5,
    "comment": "แอปพลิเคชันใช้งานง่าย สะดวก รวดเร็วดีมากค่ะ",
    "category": "product"
  }
  ```
* **คำสั่ง curl:**
  ```bash
  curl -X POST http://localhost:3000/api/feedbacks \
    -H "Content-Type: application/json" \
    -d '{
      "customer_id": "60d5ecb862b53b3b44747209",
      "rating": 5,
      "comment": "แอปพลิเคชันใช้งานง่าย สะดวก รวดเร็วดีมากค่ะ",
      "category": "product"
    }'
  ```
* **ตัวอย่าง Success Response (200 OK / 201 Created):**
  *(ระบบจะคำนวณ `sentiment` (อารมณ์/ความรู้สึก) ให้อัตโนมัติจากข้อความที่ส่งมา)*
  ```json
  {
    "success": true,
    "message": "feedback created",
    "data": {
      "id": "60d5ecd062b53b3b4474720c",
      "customer_id": "60d5ecb862b53b3b44747209",
      "rating": 5,
      "comment": "แอปพลิเคชันใช้งานง่าย สะดวก รวดเร็วดีมากค่ะ",
      "category": "product",
      "sentiment": "positive",
      "created_at": "2026-05-22T01:39:12Z"
    }
  }
  ```

---

## 📞 5. ข้อมูลการติดตามผล (Follow-Up API)

ใช้บันทึกประวัติการโทรติดตาม และอัปเดตสถานะการดูแลลูกค้า

### **5.1 บันทึกการติดตามผลใหม่ (Create Follow-Up)**
* **Method:** `POST`
* **Path:** `/api/follow-ups`
* **Headers:** `Content-Type: application/json`
* **เงื่อนไขการ Validation:**
  * `customer_id` (ต้องกรอก)
  * `type` (ต้องเป็นคำใดคำหนึ่งดังต่อไปนี้: `payment_remind`, `feedback_reply`, `general`)
  * `note` (ต้องกรอกรายละเอียดการติดตามผล)
* **โครงสร้าง Request Body (JSON):**
  ```json
  {
    "customer_id": "60d5ecb862b53b3b44747209",
    "type": "payment_remind",
    "note": "โทรแจ้งเตือนยอดค้างชำระ ลูกค้าแจ้งว่าจะชำระภายในพรุ่งนี้เช้า"
  }
  ```
* **คำสั่ง curl:**
  ```bash
  curl -X POST http://localhost:3000/api/follow-ups \
    -H "Content-Type: application/json" \
    -d '{
      "customer_id": "60d5ecb862b53b3b44747209",
      "type": "payment_remind",
      "note": "โทรแจ้งเตือนยอดค้างชำระ ลูกค้าแจ้งว่าจะชำระภายในพรุ่งนี้เช้า"
    }'
  ```
* **ตัวอย่าง Success Response (200 OK):**
  *(สถานะเริ่มต้นจะเป็น `pending` โดยอัตโนมัติ)*
  ```json
  {
    "success": true,
    "message": "follow up created",
    "data": {
      "id": "60b94326543b3b44747203fa",
      "customer_id": "60d5ecb862b53b3b44747209",
      "type": "payment_remind",
      "note": "โทรแจ้งเตือนยอดค้างชำระ ลูกค้าแจ้งว่าจะชำระภายในพรุ่งนี้เช้า",
      "status": "pending",
      "created_at": "2026-05-22T01:39:12Z"
    }
  }
  ```

---

### **5.2 อัปเดตสถานะการติดตามผล (Update Follow-Up Status)**
ใช้เพื่อเปลี่ยนสถานะของใบงานติดตามผล เช่น ดำเนินการเสร็จสิ้นแล้ว (`done`) หรือรอการติดตาม (`pending`)
* **Method:** `PATCH`
* **Path:** `/api/follow-ups/:id` (แทนที่ `:id` ด้วย ObjectID ของ Follow-up)
* **Headers:** `Content-Type: application/json`
* **เงื่อนไขการ Validation:**
  * `status` (ต้องกรอก และต้องเป็นคำใดคำหนึ่งดังต่อไปนี้: `pending`, `done`)
* **โครงสร้าง Request Body (JSON):**
  ```json
  {
    "status": "done"
  }
  ```
* **คำสั่ง curl:**
  ```bash
  curl -X PATCH http://localhost:3000/api/follow-ups/60b94326543b3b44747203fa \
    -H "Content-Type: application/json" \
    -d '{
      "status": "done"
    }'
  ```
* **ตัวอย่าง Success Response (200 OK):**
  ```json
  {
    "success": true,
    "message": "follow up status updated",
    "data": {
      "id": "60b94326543b3b44747203fa",
      "customer_id": "60d5ecb862b53b3b44747209",
      "type": "payment_remind",
      "note": "โทรแจ้งเตือนยอดค้างชำระ ลูกค้าแจ้งว่าจะชำระภายในพรุ่งนี้เช้า",
      "status": "done",
      "created_at": "2026-05-22T01:39:12Z"
    }
  }
  ```

---

## ⚠️ 6. ตัวอย่างผลลัพธ์ของข้อผิดพลาด (Common Error Responses)

เมื่อการร้องขอเกิดข้อผิดพลาด ระบบจะคืนสถานะ HTTP Status Code พร้อม JSON รูปแบบมาตรฐาน ดังนี้:

### **6.1 การส่งข้อมูลไม่ตรงเงื่อนไข (Validation Error - 400 Bad Request)**
เกิดขึ้นเมื่อ Request Body ข้อมูลไม่ครบถ้วน หรือมีรูปแบบไม่ตรงตามกำหนด
* **รูปแบบ JSON ที่ได้รับ:**
  ```json
  {
    "success": false,
    "message": "validation failed",
    "error": "ERR_VALIDATION",
    "data": [
      {
        "field": "rating",
        "message": "rating must be at most 5"
      },
      {
        "field": "category",
        "message": "category must be one of: service payment product branch"
      }
    ]
  }
  ```

### **6.2 รูปแบบ JSON ผิดพลาด (Parse Body Error - 400 Bad Request)**
เกิดขึ้นเมื่อรูปแบบของ JSON syntax เสียหาย เช่น ลืมเครื่องหมายจุลภาค `,` หรือวงเล็บปิด
* **รูปแบบ JSON ที่ได้รับ:**
  ```json
  {
    "success": false,
    "message": "invalid request body",
    "error": "ERR_PARSE"
  }
  ```

### **6.3 ค้นหาไม่พบ (Not Found Error - 404 Not Found)**
เกิดขึ้นเมื่อส่ง ID ที่ไม่มีอยู่ในฐานข้อมูล
* **รูปแบบ JSON ที่ได้รับ:**
  ```json
  {
    "success": false,
    "message": "customer not found",
    "error": "ERR_NOT_FOUND"
  }
  ```

### **6.4 ข้อผิดพลาดของเซิร์ฟเวอร์ (Server Error - 500 Internal Server Error)**
เกิดขึ้นเมื่อเกิดปัญหาในการติดต่อฐานข้อมูล หรือข้อผิดพลาดอื่นๆ ที่ระบบไม่ได้คาดคิด
* **รูปแบบ JSON ที่ได้รับ:**
  ```json
  {
    "success": false,
    "message": "failed to create feedback",
    "error": "mongo: no documents in result"
  }
  ```

---

## ⚡ 7. ตัวอย่างแบบด่วนสำหรับส่วนขยาย REST Client (VS Code)

คุณสามารถนำเนื้อหาด้านล่างนี้ไปคัดลอกบันทึกเป็นไฟล์ `.http` เพื่อใช้ส่วนขยาย **REST Client** ใน VS Code เพื่อยิงทดสอบแบบ Interactive ได้โดยตรง:

```http
@baseUrl = http://localhost:3000
@customerId = 60d5ecb862b53b3b44747209
@followUpId = 60b94326543b3b44747203fa

### 1. Health Check
GET {{baseUrl}}/health

### 2. List Customers
GET {{baseUrl}}/api/customers

### 3. List Customers with Filter
GET {{baseUrl}}/api/customers?branch=Bangkok&status=active

### 4. Get Customer Detail
GET {{baseUrl}}/api/customers/{{customerId}}

### 5. Get Dashboard Summary
GET {{baseUrl}}/api/stats/summary

### 6. Get Branch Stats
GET {{baseUrl}}/api/stats/by-branch

### 7. List Feedbacks
GET {{baseUrl}}/api/feedbacks

### 8. List Feedbacks with Filter
GET {{baseUrl}}/api/feedbacks?category=service&rating=5

### 9. Create Feedback
POST {{baseUrl}}/api/feedbacks
Content-Type: application/json

{
  "customer_id": "{{customerId}}",
  "rating": 5,
  "comment": "เจ้าหน้าที่ตอบคำถามชัดเจนดีมากครับ",
  "category": "service"
}

### 10. Create Follow-Up
POST {{baseUrl}}/api/follow-ups
Content-Type: application/json

{
  "customer_id": "{{customerId}}",
  "type": "feedback_reply",
  "note": "ส่งอีเมลยืนยันการรับเรื่องแนะนำแก่ทีมพัฒนาแล้ว"
}

### 11. Update Follow-Up Status
PATCH {{baseUrl}}/api/follow-ups/{{followUpId}}
Content-Type: application/json

{
  "status": "done"
}
```
