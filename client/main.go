package main

import (
	"context"
	"log"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	orderpb "my-grpc-project/order" // استيراد حزمة الطلبات
)

func main() {
	log.Println("Starting client to call Order Service...")

	// الاتصال بخدمة الطلبات على المنفذ الجديد 50052
	addr := "localhost:50052"
	conn, err := grpc.Dial(addr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("Failed to connect: %v", err)
	}
	defer conn.Close()

	// إنشاء عميل لخدمة الطلبات
	c := orderpb.NewOrderServiceClient(conn)

	// استدعاء دالة CreateOrder
	res, err := c.CreateOrder(context.Background(), &orderpb.CreateOrderRequest{
		UserName: "Alice",
		Item:     "Laptop",
	})
	if err != nil {
		log.Fatalf("Could not create order: %v", err)
	}

	// طباعة الرد النهائي المستلم من خدمة الطلبات
	log.Printf("Order Status: %s", res.GetStatusMessage())
}