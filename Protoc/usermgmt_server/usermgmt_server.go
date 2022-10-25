package main

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"math"
	"math/rand"
	"net"
	"sync"
	"time"

	pb "example.com/go-usermgmt-grpc/usermgmt"
	"github.com/golang/protobuf/proto"
	"github.com/keploy/go-sdk/integrations/kgrpcserver"
	"github.com/keploy/go-sdk/keploy"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
	// "google.golang.org/grpc/codes"
	// "google.golang.org/grpc/status"
)

const (
	port = ":50051"
)

type UserManagementServer struct {
	pb.UnimplementedUserManagementServer
	savedFeatures []*pb.Feature // read-only after initialized

	mu         sync.Mutex // protects routeNotes
	routeNotes map[string][]*pb.RouteNote
}

func (s *UserManagementServer) CreateNewUser(ctx context.Context, in *pb.NewUser) (*pb.User, error) {
	log.Printf("Received: %v", in.GetName())
	var id int32 = int32(rand.Intn(1000))
	return &pb.User{Name: in.GetName(), Age: in.GetAge(), Id: id}, nil
}

//Unary
func (s *UserManagementServer) GetFeature(ctx context.Context, point *pb.Point) (*pb.Feature, error) {
	for _, feature := range s.savedFeatures {
		if proto.Equal(feature.Location, point) {
			return feature, nil
		}
	}
	return &pb.Feature{Location: point}, nil
}

//Server side
func (s *UserManagementServer) ListFeatures(rect *pb.Rectangle, stream pb.UserManagement_ListFeaturesServer) error {
	for _, feature := range s.savedFeatures {
		if inRange(feature.Location, rect) {
			if err := stream.Send(feature); err != nil {
				return err
			}
		}
	}
	return nil
}

func inRange(point *pb.Point, rect *pb.Rectangle) bool {
	left := math.Min(float64(rect.Lo.Longitude), float64(rect.Hi.Longitude))
	right := math.Max(float64(rect.Lo.Longitude), float64(rect.Hi.Longitude))
	top := math.Max(float64(rect.Lo.Latitude), float64(rect.Hi.Latitude))
	bottom := math.Min(float64(rect.Lo.Latitude), float64(rect.Hi.Latitude))

	if float64(point.Longitude) >= left &&
		float64(point.Longitude) <= right &&
		float64(point.Latitude) >= bottom &&
		float64(point.Latitude) <= top {
		return true
	}
	return false
}

//Client-side
func (s *UserManagementServer) RecordRoute(stream pb.UserManagement_RecordRouteServer) error {
	var pointCount, featureCount, distance int32
	var lastPoint *pb.Point
	startTime := time.Now()
	for {
		point, err := stream.Recv()
		if err == io.EOF {
			endTime := time.Now()
			return stream.SendAndClose(&pb.RouteSummary{
				PointCount:   pointCount,
				FeatureCount: featureCount,
				Distance:     distance,
				ElapsedTime:  int32(endTime.Sub(startTime).Seconds()),
			})
		}
		if err != nil {
			return err
		}
		pointCount++
		for _, feature := range s.savedFeatures {
			if proto.Equal(feature.Location, point) {
				featureCount++
			}
		}
		if lastPoint != nil {
			distance += calcDistance(lastPoint, point)
		}
		lastPoint = point
	}
}

func toRadians(num float64) float64 {
	return num * math.Pi / float64(180)
}

func calcDistance(p1 *pb.Point, p2 *pb.Point) int32 {
	const CordFactor float64 = 1e7
	const R = float64(6371000) // earth radius in metres
	lat1 := toRadians(float64(p1.Latitude) / CordFactor)
	lat2 := toRadians(float64(p2.Latitude) / CordFactor)
	lng1 := toRadians(float64(p1.Longitude) / CordFactor)
	lng2 := toRadians(float64(p2.Longitude) / CordFactor)
	dlat := lat2 - lat1
	dlng := lng2 - lng1

	a := math.Sin(dlat/2)*math.Sin(dlat/2) +
		math.Cos(lat1)*math.Cos(lat2)*
			math.Sin(dlng/2)*math.Sin(dlng/2)
	c := 2 * math.Atan2(math.Sqrt(a), math.Sqrt(1-a))

	distance := R * c
	return int32(distance)
}

//Bidiorectional
func (s *UserManagementServer) RouteChat(stream pb.UserManagement_RouteChatServer) error {
	for {
		in, err := stream.Recv()
		if err == io.EOF {
			return nil
		}
		if err != nil {
			return err
		}
		key := serialize(in.Location)
		s.mu.Lock()
		s.routeNotes[key] = append(s.routeNotes[key], in)
		// Note: this copy prevents blocking other clients while serving this one.
		// We don't need to do a deep copy, because elements in the slice are
		// insert-only and never modified.
		rn := make([]*pb.RouteNote, len(s.routeNotes[key]))
		copy(rn, s.routeNotes[key])
		s.mu.Unlock()

		for _, note := range rn {
			if err := stream.Send(note); err != nil {
				return err
			}
		}
	}
}

func serialize(point *pb.Point) string {
	return fmt.Sprintf("%d %d", point.Latitude, point.Longitude)
}

func main() {
	k := keploy.New(keploy.Config{
		App: keploy.AppConfig{
			Name: "grpc-server-app",
			Port: port,
		},
		Server: keploy.ServerConfig{
			URL: "http://localhost:6789/api",
		},
	})
	lis, err := net.Listen("tcp", port)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	s := grpc.NewServer(kgrpcserver.UnaryInterceptor(k))
	reflection.Register(s)
	userObj := UserManagementServer{routeNotes: make(map[string][]*pb.RouteNote)}
	pb.RegisterUserManagementServer(s, &userObj)
	// var data []byte
	data := exampleData
	if err := json.Unmarshal(data, &userObj.savedFeatures); err != nil {
		log.Fatalf("Failed to load default features: %v", err)
	}
	log.Printf("server listening at : %v", lis.Addr())
	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}

var exampleData = []byte(`[{
    "location": {
        "latitude": 407838351,
        "longitude": -746143763
    },
    "name": "Patriots Path, Mendham, NJ 07945, USA"
}, {
    "location": {
        "latitude": 408122808,
        "longitude": -743999179
    },
    "name": "101 New Jersey 10, Whippany, NJ 07981, USA"
}, {
    "location": {
        "latitude": 413628156,
        "longitude": -749015468
    },
    "name": "U.S. 6, Shohola, PA 18458, USA"
}, {
    "location": {
        "latitude": 419999544,
        "longitude": -740371136
    },
    "name": "5 Conners Road, Kingston, NY 12401, USA"
}, {
    "location": {
        "latitude": 414008389,
        "longitude": -743951297
    },
    "name": "Mid Hudson Psychiatric Center, New Hampton, NY 10958, USA"
}, {
    "location": {
        "latitude": 419611318,
        "longitude": -746524769
    },
    "name": "287 Flugertown Road, Livingston Manor, NY 12758, USA"
}, {
    "location": {
        "latitude": 406109563,
        "longitude": -742186778
    },
    "name": "4001 Tremley Point Road, Linden, NJ 07036, USA"
}, {
    "location": {
        "latitude": 416802456,
        "longitude": -742370183
    },
    "name": "352 South Mountain Road, Wallkill, NY 12589, USA"
}, {
    "location": {
        "latitude": 412950425,
        "longitude": -741077389
    },
    "name": "Bailey Turn Road, Harriman, NY 10926, USA"
}, {
    "location": {
        "latitude": 412144655,
        "longitude": -743949739
    },
    "name": "193-199 Wawayanda Road, Hewitt, NJ 07421, USA"
}]`)
