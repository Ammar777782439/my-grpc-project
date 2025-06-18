package main

import (
	"context"
	"fmt"
	"log"
	"net"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	// استيراد الحزم التي تم إنشاؤها لكلتا الخدمتين
	greetpb "my-grpc-project/proto"  // خدمة الترحيب
	orderpb "my-grpc-project/order/orderpb" // خدمة الطلبات
)

// تعريف struct لخادم الطلبات
type orderServer struct {
	orderpb.UnimplementedOrderServiceServer
}

// تنفيذ دالة CreateOrder
func (s *orderServer) CreateOrder(ctx context.Context, req *orderpb.CreateOrderRequest) (*orderpb.CreateOrderReply, error) {
	log.Printf("Received CreateOrder RPC for user: %s, item: %s", req.GetUserName(), req.GetItem())
	
	// --- الجزء الأول: التصرف كـ "عميل" لخدمة الترحيب ---

	// الاتصال بخادم الترحيب Greeter Service على منفذ 50051
	conn, err := grpc.Dial("localhost:50051", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Printf("Could not connect to Greeter service: %v", err)
		// في تطبيق حقيقي، يجب معالجة هذا الخطأ بشكل أفضل
		return &orderpb.CreateOrderReply{StatusMessage: "Error: Could not process order."}, nil
	}
	defer conn.Close()

	// إنشاء عميل لخدمة الترحيب
	greeterClient := greetpb.NewGreeterClient(conn)

	// استدعاء دالة SayHello من خدمة الترحيب
	greetRes, err := greeterClient.SayHello(context.Background(), &greetpb.HelloRequest{Name: req.GetUserName()})
	if err != nil {
		log.Printf("Could not get greeting: %v", err)
		return &orderpb.CreateOrderReply{StatusMessage: "Error: Could not process order."}, nil
	}
	
	// --- الجزء الثاني: إكمال منطق خدمة الطلبات ---

	// دمج رسالة الترحيب مع تفاصيل الطلب
	finalMessage := fmt.Sprintf("%s. Your order for a '%s' has been created.", greetRes.GetMessage(), req.GetItem())

	// إرجاع الرد النهائي
	return &orderpb.CreateOrderReply{StatusMessage: finalMessage}, nil
}


func main() {
	// بدء خادم الطلبات على منفذ 50052
	lis, err := net.Listen("tcp", ":50052")
	if err != nil {
		log.Fatalf("Failed to listen on port 50052: %v", err)
	}

	s := grpc.NewServer()
	orderpb.RegisterOrderServiceServer(s, &orderServer{})

	log.Println("Order Service listening on port 50052")
	if err := s.Serve(lis); err != nil {
		log.Fatalf("Failed to serve Order Service: %v", err)
	}
}