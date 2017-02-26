package shared 

 
import (
 // "vpms"
)

 type Args struct { 
 	A, B int 
 } 
 
 
 type Quotient struct { 
 	Quo, Rem int 
 }

type Pargs struct{
  Name string
 }

 type Pinputargs struct{
  Name  string
  Value string
 }
 type Pcalcargs struct{
  Calc string
 }
 type Resp struct{
   Result string
   Message string
   Field string
 }

 type Arith interface { 
     Multiply(args *Args, reply *int) error 
     Divide(args *Args, quo *Quotient) error 
 } 
 
type Prod interface { 
     Init( args *Pargs , reply int) error
     AddInput( args *Pinputargs , reply *int) error
     Calc( args *Pcalcargs , reply *Resp ) error
 } 

