package main

import (
	// "context" يوفر وسيلة لإرسال إشارات الإلغاء والمواعيد النهائية عبر الطلبات.
	"context" 
	// "log" لطباعة الرسائل على الشاشة (للتصحيح والمراقبة).
	"log"
	// "net" يوفر واجهة للشبكات، سنستخدمه لفتح منفذ TCP.
	"net"

	// هذه هي مكتبة gRPC الرئيسية التي قمنا بتثبيتها.
	"google.golang.org/grpc"

	// هذا هو الكود الذي تم إنشاؤه من ملف .proto الخاص بنا.
	// المسار هو "اسم-المشروع/المجلد/الحزمة".
	// اسم الحزمة (greetpb) حددناه في خيار go_package داخل ملف .proto.
	greetpb "my-grpc-project/proto"
)

// نعرّف struct اسمه 'server'. هذا الـ struct سيمثل خادمنا المنطقي.
// سنقوم بتضمين 'UnimplementedGreeterServer' لضمان التوافق مع الإصدارات المستقبلية.
// إذا أضفت وظيفة جديدة للخدمة في ملف .proto ولم تنفذها هنا، لن يفشل الكود عند الترجمة.
type server struct {
	greetpb.UnimplementedGreeterServer
}

// هنا نقوم بتنفيذ (implement) دالة SayHello التي عرفناها في ملف .proto.
// هذه هي الوظيفة الحقيقية التي سيتم استدعاؤها عندما يرسل العميل (Client) طلباً.
// تستقبل 'context' و 'HelloRequest' كمدخلات.
func (s *server) SayHello(ctx context.Context, req *greetpb.HelloRequest) (*greetpb.HelloReply, error) {
	// نطبع رسالة على الخادم لنعرف أن الدالة قد تم استدعاؤها.
	log.Printf("Received SayHello RPC for name: %v", req.GetName())

	// نستخرج الاسم من الطلب القادم.
	name := req.GetName()
	// ننشئ رسالة الترحيب.
	message := "Hello " + name

	// نرجع رسالة من نوع HelloReply تحتوي على الترحيب، ولا نرجع أي خطأ (nil).
	// يجب أن يتطابق نوع الإرجاع مع ما حددناه في ملف .proto.
	return &greetpb.HelloReply{Message: message}, nil
}

func main() {
	// نطبع رسالة لبدء تشغيل الخادم.
	log.Println("Starting gRPC server...")

	// نستمع على منفذ TCP رقم 50051. هذا هو العنوان الذي سينتظر الخادم عليه الطلبات.
	// "tcp" هو البروتوكول، و ":50051" يعني كل الواجهات على المنفذ 50051.
	lis, err := net.Listen("tcp", ":50051")
	if err != nil {
		// إذا فشل فتح المنفذ (لأنه مستخدم مثلاً)، نوقف البرنامج فوراً.
		log.Fatalf("Failed to listen: %v", err)
	}

	// ننشئ نسخة جديدة من خادم gRPC.
	s := grpc.NewServer()

	// نسجل خدمتنا (server struct) مع خادم gRPC.
	// هذا السطر يخبر خادم gRPC: "عندما يأتي طلب لخدمة Greeter، استخدم هذا الـ server لتنفيذه".
	greetpb.RegisterGreeterServer(s, &server{})

	// نطبع عنوان الخادم الذي يستمع عليه.
	log.Printf("Server listening at %v", lis.Addr())

	// نبدأ الخادم لكي يستقبل الطلبات على المنفذ الذي فتحناه.
	// s.Serve() هي عملية حاصرة (blocking)، أي أن البرنامج سيبقى في هذا السطر إلى الأبد
	// يستمع للطلبات، إلا إذا حدث خطأ.
	if err := s.Serve(lis); err != nil {
		log.Fatalf("Failed to serve: %v", err)
	}
}