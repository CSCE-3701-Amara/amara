// Command demo-seed creates synthetic, isolated company data for demonstrations.
package main

import (
	"context"
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	_ "embed"
	"encoding/base32"
	"encoding/base64"
	"encoding/json"
	"flag"
	"fmt"
	"net/url"
	"os"

	"github.com/CSCE-3701-Amara/amara/internal/invw/domain"
	"github.com/jackc/pgx/v5"
	"golang.org/x/crypto/argon2"
)

func run() error {
	apply := flag.Bool("apply", false, "commit the demo data (default rolls back)")
	activate := flag.Bool("activate", false, "validate existing synthetic modules and activate four demo pipelines")
	parallel := flag.Bool("parallel", false, "add two published parallel-branch demo jobs")
	host := flag.String("allow-host", "", "exact database hostname explicitly approved for seeding")
	flag.Parse()
	address := os.Getenv("DATABASE_URL")
	u, err := url.Parse(address)
	if err != nil || u.Hostname() == "" || *host == "" || u.Hostname() != *host {
		return fmt.Errorf("DATABASE_URL hostname must exactly match --allow-host")
	}
	passwordBytes := make([]byte, 18)
	salt := make([]byte, 16)
	if _, err = rand.Read(passwordBytes); err != nil {
		return err
	}
	if _, err = rand.Read(salt); err != nil {
		return err
	}
	password := base64.RawURLEncoding.EncodeToString(passwordBytes)
	hash := argon2.IDKey([]byte(password), salt, 3, 64*1024, 1, 32)
	mfaKey := []byte(os.Getenv("AUTH_MFA_ENCRYPTION_KEY"))
	if len(mfaKey) != 32 {
		return fmt.Errorf("AUTH_MFA_ENCRYPTION_KEY must be exactly 32 bytes")
	}
	secret := make([]byte, 20)
	if _, err = rand.Read(secret); err != nil {
		return err
	}
	block, err := aes.NewCipher(mfaKey)
	if err != nil {
		return err
	}
	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return err
	}
	nonce := make([]byte, gcm.NonceSize())
	if _, err = rand.Read(nonce); err != nil {
		return err
	}
	encryptedSecret := gcm.Seal(nonce, nonce, secret, nil)
	ctx := context.Background()
	config, err := pgx.ParseConfig(address)
	if err != nil {
		return fmt.Errorf("parse database configuration: %w", err)
	}
	config.DefaultQueryExecMode = pgx.QueryExecModeSimpleProtocol
	conn, err := pgx.ConnectConfig(ctx, config)
	if err != nil {
		return fmt.Errorf("connect database: %w", err)
	}
	defer conn.Close(ctx)
	tx, err := conn.Begin(ctx)
	if err != nil {
		return err
	}
	defer tx.Rollback(ctx)
	var exists bool
	if err = tx.QueryRow(ctx, "SELECT EXISTS (SELECT 1 FROM companies WHERE lower(email) = 'demo.company@amara.example.com')").Scan(&exists); err != nil {
		return err
	}
	if exists {
		if *parallel {
			if _, err = tx.Exec(ctx, parallelSQL); err != nil {
				return fmt.Errorf("parallel demo: %w", err)
			}
			if *apply {
				if err = tx.Commit(ctx); err != nil {
					return err
				}
				fmt.Println("Parallel demo jobs ready; existing credentials preserved.")
			} else {
				fmt.Println("Parallel demo dry run passed; transaction rolled back.")
			}
			return nil
		}
		if *activate {
			if _, err = tx.Exec(ctx, activateSQL); err != nil {
				return err
			}
			rows, queryErr := tx.Query(ctx, "SELECT authoring_definition FROM modules WHERE company_id=md5('amara-demo-company-v1')::uuid")
			if queryErr != nil {
				return queryErr
			}
			for rows.Next() {
				var raw []byte
				if err = rows.Scan(&raw); err != nil {
					rows.Close()
					return err
				}
				var definition domain.Definition
				if err = json.Unmarshal(raw, &definition); err != nil {
					rows.Close()
					return err
				}
				if err = definition.ValidatePublication(nil); err != nil {
					rows.Close()
					return fmt.Errorf("module %s: %w", definition.Name, err)
				}
			}
			rows.Close()
			if err = rows.Err(); err != nil {
				return err
			}
			if *apply {
				if err = tx.Commit(ctx); err != nil {
					return err
				}
				fmt.Println("Four demo pipelines activated; original demo credentials preserved.")
			} else {
				fmt.Println("Activation dry run passed; transaction rolled back.")
			}
			return nil
		}
		return fmt.Errorf("demo company already exists; no records were changed")
	}
	if _, err = tx.Exec(ctx, seedSQL, hash, salt, encryptedSecret); err != nil {
		return fmt.Errorf("seed transaction: %w", err)
	}
	if !*apply {
		fmt.Println("Dry run passed; transaction rolled back.")
		return nil
	}
	if err = tx.Commit(ctx); err != nil {
		return err
	}
	fmt.Printf("Created demo company and candidates. Demo password: %s\n", password)
	fmt.Printf("Company authenticator setup key: %s\n", base32.StdEncoding.WithPadding(base32.NoPadding).EncodeToString(secret))
	fmt.Println("Company: demo.company@amara.example.com; candidates: demo.candidate01@amara.example.com through demo.candidate24@amara.example.com")
	return nil
}

//go:embed parallel.sql
var parallelSQL string

const activateSQL = `
UPDATE modules SET authoring_definition=authoring_definition-'rubric'
WHERE id IN (md5('amara-demo-module-code-v1')::uuid,md5('amara-demo-module-behavior-v1')::uuid);
UPDATE modules SET authoring_definition=authoring_definition-'amarat'
WHERE id=md5('amara-demo-module-behavior-v1')::uuid;
UPDATE test_configurations tc SET configuration=jsonb_build_object('position',0,'definition',m.authoring_definition)
FROM modules m WHERE tc.test_module_id=m.id AND tc.stage_id IN (
SELECT id FROM pipeline_stages WHERE version_id IN (SELECT md5('amara-demo-version-v1-'||n)::uuid FROM generate_series(1,5) n));
UPDATE pipeline_stages SET scheduled_start_at=now()-interval '1 day',scheduled_end_at=now()+interval '30 days'
WHERE version_id IN (SELECT md5('amara-demo-version-v1-'||n)::uuid FROM generate_series(1,4) n);
UPDATE pipeline_versions SET is_published=true,published_at=now(),published_by=md5('amara-demo-company-v1')::uuid
WHERE id IN (SELECT md5('amara-demo-version-v1-'||n)::uuid FROM generate_series(1,4) n);
UPDATE pipelines SET published_version_id=md5('amara-demo-version-v1-'||n)::uuid,status='published'
FROM generate_series(1,4) n WHERE id=md5('amara-demo-pipeline-v1-'||n)::uuid;
`

func main() {
	if err := run(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

const seedSQL = `
INSERT INTO companies (id,email,password_hash,salt,legal_name,display_name,description,approved,created_at)
VALUES (md5('amara-demo-company-v1')::uuid,'demo.company@amara.example.com',$1,$2,'Nile Horizon Demo Ltd','Nile Horizon Demo','Synthetic company for hiring demonstrations.',true,now()-interval '120 days');
INSERT INTO auth_mfa_credentials(account_id,account_type,secret,enabled,enabled_at)
VALUES (md5('amara-demo-company-v1')::uuid,'company',$3,true,now());
INSERT INTO candidates (id,name,email,password_hash,salt,created_at,account_status,nationality,address,biography,years_of_experience,available)
SELECT md5('amara-demo-candidate-v1-'||n)::uuid,
 (ARRAY['Mariam','Omar','Salma','Youssef','Nour','Karim','Laila','Hassan','Farah','Ahmed','Rana','Tarek'])[(n-1)%12+1] || ' ' || (ARRAY['Adel','Hassan','Mostafa','Sami'])[(n-1)/12+1],
 'demo.candidate'||lpad(n::text,2,'0')||'@amara.example.com',$1,$2,now()-make_interval(days=>100+n),'active','Egyptian',
 (ARRAY['Cairo','Alexandria','Giza','Mansoura'])[(n-1)%4+1],
 'Synthetic candidate interested in building practical, reliable products.', (n-1)%9,true
FROM generate_series(1,24) n;
INSERT INTO job_postings(id,company_id,title,description,language,currency,seniority_level,employment_type,location,salary_min,salary_max,salary_visible,open_positions,application_opening,application_deadline,max_applications,finalist_target,status,published_at,created_at,updated_at)
SELECT md5('amara-demo-job-v1-'||n)::uuid,md5('amara-demo-company-v1')::uuid,
 (ARRAY['Backend Engineer','Frontend Engineer','Data Analyst','Product Designer','QA Engineer','DevOps Engineer'])[(n-1)%6+1],
 'Join a collaborative team solving real customer problems. This is a synthetic demonstration posting.', 'English','EGP',
 (ARRAY['Junior','Mid-level','Senior'])[(n-1)%3+1], 'remote full-time','Cairo',20000+n*1000,35000+n*1500,true,2,
 now()-interval '14 days',now()+interval '45 days',100,3,CASE WHEN n=6 THEN 'draft' ELSE 'published' END,
 CASE WHEN n=6 THEN NULL ELSE now()-interval '14 days' END,now()-interval '20 days',now()-interval '1 day'
FROM generate_series(1,6) n;
INSERT INTO modules(id,company_id,name,description,type,published,authoring_definition)
VALUES (md5('amara-demo-module-code-v1')::uuid,md5('amara-demo-company-v1')::uuid,'Practical coding','A short coding exercise.','coding',true,
 '{"name":"Practical coding","description":"A short coding exercise.","type":"coding","active_minutes":30,"threshold":70,"questions":[{"prompt":"Read two integers and print their sum.","format":"code","difficulty":"easy","amara":"Problem solving","test_cases":[{"input":"2 3","expected_output":"5","weight":1},{"input":"-1 1","expected_output":"0","weight":1}]}],"amarat":[{"name":"Problem solving","crest":"logic","domain":"Technical","description":"Solve practical problems"}],"rubric":[{"description":"Produces correct results","amara":"Problem solving"}]}'::jsonb),
 (md5('amara-demo-module-behavior-v1')::uuid,md5('amara-demo-company-v1')::uuid,'Collaboration and communication','Describe how you work with others.','behavioral',true,
 '{"name":"Collaboration and communication","description":"Describe how you work with others.","type":"behavioral","active_minutes":20,"questions":[{"prompt":"Describe a time you resolved a disagreement with a teammate.","format":"long_text","amara":"Collaboration"}],"amarat":[{"name":"Collaboration","crest":"team","domain":"Communication","description":"Work constructively with others"}],"rubric":[{"description":"Explains a constructive approach","amara":"Collaboration"}]}'::jsonb);
INSERT INTO pipelines(id,posting_id,status) SELECT md5('amara-demo-pipeline-v1-'||n)::uuid,md5('amara-demo-job-v1-'||n)::uuid,'draft' FROM generate_series(1,5) n;
INSERT INTO pipeline_versions(id,pipeline_id,version_number,is_published) SELECT md5('amara-demo-version-v1-'||n)::uuid,md5('amara-demo-pipeline-v1-'||n)::uuid,1,false FROM generate_series(1,5) n;
UPDATE job_postings SET pipeline_id=md5('amara-demo-pipeline-v1-'||n)::uuid FROM generate_series(1,5) n WHERE id=md5('amara-demo-job-v1-'||n)::uuid;
INSERT INTO pipeline_stages(id,version_id,name,module_type,advancement_rule,stage_order,position_x,position_y)
SELECT md5('amara-demo-stage-v1-'||n||'-'||s)::uuid,md5('amara-demo-version-v1-'||n)::uuid,
 (ARRAY['Coding assessment','Collaboration assessment','Hiring manager review']) [s], 'assessment',CASE WHEN s=3 THEN 'manual' ELSE 'completion' END,s,40+(s-1)*260,40
FROM generate_series(1,5) n CROSS JOIN generate_series(1,3) s;
INSERT INTO stage_transitions(id,pipeline_version_id,from_stage_id,to_stage_id)
SELECT md5('amara-demo-edge-v1-'||n||'-'||s)::uuid,md5('amara-demo-version-v1-'||n)::uuid,md5('amara-demo-stage-v1-'||n||'-'||s)::uuid,md5('amara-demo-stage-v1-'||n||'-'||(s+1))::uuid
FROM generate_series(1,5) n CROSS JOIN generate_series(1,2) s;
INSERT INTO test_configurations(id,stage_id,test_module_id,active_time_limit,configuration)
SELECT md5('amara-demo-config-v1-'||n||'-'||s)::uuid,md5('amara-demo-stage-v1-'||n||'-'||s)::uuid,
 CASE WHEN s=1 THEN md5('amara-demo-module-code-v1')::uuid ELSE md5('amara-demo-module-behavior-v1')::uuid END,
 interval '30 minutes',jsonb_build_object('position',0,'definition',(SELECT authoring_definition FROM modules WHERE id=CASE WHEN s=1 THEN md5('amara-demo-module-code-v1')::uuid ELSE md5('amara-demo-module-behavior-v1')::uuid END))
FROM generate_series(1,5) n CROSS JOIN generate_series(1,3) s;
INSERT INTO applications(id,posting_id,cand_id,status,created_at,updated_at,submitted_at,screened_out,screening_completed_at)
SELECT md5('amara-demo-application-v1-'||n)::uuid,md5('amara-demo-job-v1-'||((n-1)%5+1))::uuid,md5('amara-demo-candidate-v1-'||n)::uuid,
 'admitted',now()-interval '5 days',now()-interval '5 days',now()-interval '5 days',false,now()-interval '5 days'
FROM generate_series(1,24) n;
`
