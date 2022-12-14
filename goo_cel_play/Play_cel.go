package goo_cel_play

import (
	"encoding/json"
	"fmt"
	"log"
	"reflect"
	"sort"
	"strings"
	"time"

	"github.com/google/cel-go/cel"
	"github.com/google/cel-go/common/types"
	"github.com/google/cel-go/common/types/ref"
	"google.golang.org/protobuf/encoding/protojson"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/types/known/structpb"
)

// codelab: https://codelabs.developers.google.com/codelabs/cel-go

func PlayCel() {
	fmt.Print("============================================\n \t play1 - hello \n============================================\n\n")
	test1_simple(`"Hello Go"`)
	fmt.Print("============================================\n \t play2 - with variable + proto struct \n============================================\n\n")
	test2_variable()
	fmt.Print("============================================\n \t play3 - with func \n============================================\n\n")
	test3_customFunc()
	fmt.Print("============================================\n \t play4 - with json gen \n============================================\n\n")
	test4_jsonGen()
}

func test4_jsonGen() {

	env, err := cel.NewEnv(
		cel.Variable("now", cel.TimestampType),
	)
	if err != nil {
		log.Panic(err.Error())
	}

	ast := compile(env, `
		{'sub': 'serviceAccount:delegate@acme.co',
		'aud': 'my-project',
		'iss': 'auth.acme.com:12350',
		'iat': now,
		'nbf': now,
		'exp': now + duration('300s'),
		'extra_claims': {
				'group': 'admin'
		}}
	`, cel.MapType(cel.StringType, cel.DynType))
	program, _ := env.Program(ast)

	dumpAts2Json(ast)

	// Evaluate a request object that sets the proper group claim.
	out, _, _ := eval(program, map[string]any{
		"now": time.Now(),
	})

	log.Println(valueToJSON(out))
}

func test3_customFunc() {
	// type-signature for 'endWiths'.
	typeParamA := cel.TypeParamType("A")

	custom_str_endWiths := func(args ...ref.Val) ref.Val {
		this := strings.ToLower(string(args[0].(types.String)))
		inp := strings.ToLower(string(args[1].(types.String)))

		fmt.Printf("custom_str_endWiths: this: %+v inp: %+v\n", this, inp)

		if strings.HasSuffix(this, inp) {
			return types.True
		}

		return types.False
	}

	env, err := cel.NewEnv(
		cel.Types(&RequestContext{}),
		cel.Variable("request",
			cel.ObjectType("play2022.goo_cel_play.requestContext"),
		),
		cel.Function("endWiths",
			cel.MemberOverload(
				"string_endWiths",
				[]*cel.Type{cel.StringType, typeParamA},
				cel.BoolType,
				cel.FunctionBinding(custom_str_endWiths),
			),
		),
	)

	if err != nil {
		log.Fatal(err)
	}
	ast := compile(env, `request.email.endWiths('@acm.com')
	`, cel.BoolType)
	program, _ := env.Program(ast)

	dumpAts2Json(ast)

	// Evaluate a request object that sets the proper group claim.
	eval(program, map[string]any{
		"request": &RequestContext{
			Email: "abc@acm.com",
		},
	})

	eval(program, map[string]any{
		"request": &RequestContext{
			Email: "abc@dummy.com",
		},
	})

}

func test2_variable() {

	env, err := cel.NewEnv(
		cel.Types(&RequestContext{}),
		cel.Variable("request",
			cel.ObjectType("play2022.goo_cel_play.requestContext"),
		),
	)
	if err != nil {
		log.Fatal(err)
	}
	ast := compile(env, `request.group == 'admin'
	|| request.email == 'super@admin.universal'
	`, cel.BoolType)

	dumpAts2Json(ast)

	program, _ := env.Program(ast)

	// Evaluate a request object that sets the proper group claim.
	eval(program, map[string]any{
		"request": &RequestContext{
			Group: "user",
		},
	})

	eval(program, map[string]any{
		"request": &RequestContext{
			Group: "admin",
		},
	})
	eval(program, map[string]any{
		"request": &RequestContext{
			Group: "unknown",
			Email: "super@admin.universal",
		},
	})
}

func test1_simple(exp string) {
	// Create the standard environment.
	env, err := cel.NewEnv()
	if err != nil {
		log.Panic(err)
	}
	ast, output, detail := doEvalAndRunExpWithEnv(exp, env)

	dumpAts2Json(ast)

	log.Printf("Res: %+v", map[string]any{
		"output": output,
		"detail": detail,
	})
}

func doEvalAndRunExpWithEnv(exp string, env *cel.Env) (*cel.Ast, ref.Val, *cel.EvalDetails) {
	// Parse -> ast
	ast, issue := env.Parse(exp)
	if issue != nil {
		log.Panic(issue)
	}
	log.Printf("Exp: %s\n", exp)
	log.Printf("\tAST: %+v\n", ast)
	log.Printf("\tOutputType 1: %+v", ast.OutputType())

	// ast -> typeCheck(infer type)
	ast, iss := env.Check(ast)
	if iss != nil {
		log.Panic(issue)
	}
	log.Printf("\tOutputType 2: %+v", ast.OutputType())

	// ast_checked -> run
	log.Printf("Run: %s\n", exp)
	program, err := env.Program(ast)
	if err != nil {
		log.Panic(err)
	}
	output, detail, err := program.Eval(cel.NoVars())
	if err != nil {
		log.Panic(err)
	}

	return ast, output, detail
}

// -----------------------------------------------------------------------------------------------------
// helper
// -----------------------------------------------------------------------------------------------------

func compile(env *cel.Env, expr string, celType *cel.Type) *cel.Ast {

	// Note: Do both parse and check at one call !
	ast, iss := env.Compile(expr)
	if iss.Err() != nil {
		log.Fatal(iss.Err())
	}

	if !reflect.DeepEqual(ast.OutputType(), celType) {
		log.Fatalf(
			"Got %v, wanted %v result type", ast.OutputType(), celType)
	}
	fmt.Printf("EXP: %s\n\n", strings.ReplaceAll(expr, "\t", " "))
	return ast
}

// eval will evaluate a given program `prg` against a set of variables `vars`
// and return the output, eval details (optional), or error that results from
// evaluation.
func eval(prg cel.Program,
	vars any) (out ref.Val, det *cel.EvalDetails, err error) {
	varMap, isMap := vars.(map[string]any)
	fmt.Println("------ input ------")
	if !isMap {
		fmt.Printf("(%T)\n", vars)
	} else {
		for k, v := range varMap {
			switch v.(type) {
			case map[string]any:
				b, _ := json.MarshalIndent(v, "", "  ")
				fmt.Printf("%s = %v\n", k, string(b))
			case uint64:
				fmt.Printf("%s = %vu\n", k, v)
			default:
				fmt.Printf("%s = %v\n", k, v)
			}
		}
	}
	fmt.Println()
	out, det, err = prg.Eval(vars)
	report(out, det, err)
	fmt.Println()
	return
}

// report prints out the result of evaluation in human-friendly terms.
func report(result ref.Val, details *cel.EvalDetails, err error) {
	fmt.Println("------ result ------")
	if err != nil {
		fmt.Printf("error: %s\n", err)
	} else {
		fmt.Printf("value: %v (%T)\n", result, result)
	}
	if details != nil {
		fmt.Printf("\n------ eval states ------\n")
		state := details.State()
		stateIDs := state.IDs()
		ids := make([]int, len(stateIDs), len(stateIDs))
		for i, id := range stateIDs {
			ids[i] = int(id)
		}
		sort.Ints(ids)
		for _, id := range ids {
			v, found := state.Value(int64(id))
			if !found {
				continue
			}
			fmt.Printf("%d: %v (%T)\n", id, v, v)
		}
	}
}

func dumpAts2Json(ats *cel.Ast) {
	bytes, err := protojson.Marshal(ats.Expr())
	if err != nil {
		log.Panic(err)
	}

	fmt.Printf("JSON: %s\n", string(bytes))
}

// valueToJSON converts the CEL type to a protobuf JSON representation and
// marshals the result to a string.
func valueToJSON(val ref.Val) string {
	v, err := val.ConvertToNative(reflect.TypeOf(&structpb.Value{}))
	if err != nil {
		log.Panic(err)
	}
	marshaller := protojson.MarshalOptions{Indent: "    "}
	bytes, err := marshaller.Marshal(v.(proto.Message))
	if err != nil {
		log.Panic(err)
	}
	return string(bytes)
}
