package seed

import (
	"database/sql"
	"time"

	"github.com/bedrock/services/demo/inertia/api/internal/database"
)

func seedOrganizations(db *sql.DB) ([]int64, error) {
	orgs := []string{
		"Acme Ventures",
		"Northstar Logistics",
		"Juniper Labs",
		"Summit Advisory",
		"Atlas Health",
		"Meridian Finance",
		"Cascade Networks",
		"Silverline Media",
		"Helio Systems",
		"Polar Dynamics",
		"Cedar Analytics",
		"Orchid Capital",
		"Bluebird Software",
		"Mistral Consulting",
		"Beacon Industries",
	}

	ids := make([]int64, 0, len(orgs))

	for _, org := range orgs {
		id, err := database.CreateOrganization(db, org)

		if err != nil {
			return nil, err
		}

		ids = append(ids, id)
	}

	return ids, nil
}

func seedContacts(db *sql.DB, orgIDs []int64) error {
	type contactSeed struct {
		firstName string
		lastName  string
		email     string
		phone     string
		orgID     *int64
		favorite  bool
	}

	org := func(index int) *int64 { return &orgIDs[index] }

	contacts := []contactSeed{
		// Acme Ventures
		{"Alicia", "Keys", "alicia@acme.test", "+1 555 0100", org(0), true},
		{"Brandon", "Lee", "brandon@acme.test", "+1 555 0106", org(0), false},
		{"Carmen", "Diaz", "carmen@acme.test", "+1 555 0107", org(0), false},
		{"David", "Chen", "david@acme.test", "+1 555 0108", org(0), true},
		{"Elena", "Rossi", "elena@acme.test", "+1 555 0109", org(0), false},
		// Northstar Logistics
		{"Marcus", "Tan", "marcus@northstar.test", "+1 555 0101", org(1), false},
		{"Fatima", "Al-Rashid", "fatima@northstar.test", "+1 555 0110", org(1), true},
		{"George", "Nakamura", "george@northstar.test", "+1 555 0111", org(1), false},
		{"Hannah", "Johansson", "hannah@northstar.test", "+1 555 0112", org(1), false},
		{"Isaac", "Mbeki", "isaac@northstar.test", "+1 555 0113", org(1), false},
		{"Jana", "Kowalski", "jana@northstar.test", "+1 555 0114", org(1), false},
		// Juniper Labs
		{"Priya", "Singh", "priya@juniper.test", "+1 555 0102", org(2), true},
		{"Kai", "Tanaka", "kai@juniper.test", "+1 555 0115", org(2), false},
		{"Lena", "Vogt", "lena@juniper.test", "+1 555 0116", org(2), false},
		{"Mateo", "Garcia", "mateo@juniper.test", "+1 555 0117", org(2), false},
		{"Nina", "Petrov", "nina@juniper.test", "+1 555 0118", org(2), true},
		{"Oscar", "Fernandez", "oscar@juniper.test", "+1 555 0119", org(2), false},
		{"Paloma", "Reyes", "paloma@juniper.test", "+1 555 0120", org(2), false},
		// Summit Advisory
		{"Jules", "Martin", "jules@summit.test", "+1 555 0103", org(3), false},
		{"Quinn", "O'Brien", "quinn@summit.test", "+1 555 0121", org(3), false},
		{"Remi", "Dubois", "remi@summit.test", "+1 555 0122", org(3), true},
		{"Sana", "Okafor", "sana@summit.test", "+1 555 0123", org(3), false},
		{"Tomás", "Silva", "tomas@summit.test", "+1 555 0124", org(3), false},
		{"Uma", "Patel", "uma@summit.test", "+1 555 0125", org(3), false},
		// Atlas Health
		{"Nora", "Alvarez", "nora@atlas.test", "+1 555 0104", org(4), false},
		{"Victor", "Chang", "victor@atlas.test", "+1 555 0126", org(4), false},
		{"Wendy", "Kim", "wendy@atlas.test", "+1 555 0127", org(4), true},
		{"Xavier", "Lopez", "xavier@atlas.test", "+1 555 0128", org(4), false},
		{"Yara", "Hassan", "yara@atlas.test", "+1 555 0129", org(4), false},
		{"Zane", "Mitchell", "zane@atlas.test", "+1 555 0130", org(4), false},
		// Meridian Finance
		{"Aiden", "Brooks", "aiden@meridian.test", "+1 555 0131", org(5), false},
		{"Bella", "Torres", "bella@meridian.test", "+1 555 0132", org(5), true},
		{"Caleb", "Wright", "caleb@meridian.test", "+1 555 0133", org(5), false},
		{"Dina", "Nguyen", "dina@meridian.test", "+1 555 0134", org(5), false},
		{"Ethan", "Osei", "ethan@meridian.test", "+1 555 0135", org(5), false},
		{"Freya", "Anderson", "freya@meridian.test", "+1 555 0136", org(5), false},
		// Cascade Networks
		{"Gael", "Moreno", "gael@cascade.test", "+1 555 0137", org(6), false},
		{"Hana", "Sato", "hana@cascade.test", "+1 555 0138", org(6), true},
		{"Ivan", "Kozlov", "ivan@cascade.test", "+1 555 0139", org(6), false},
		{"Jade", "Campbell", "jade@cascade.test", "+1 555 0140", org(6), false},
		{"Kiran", "Sharma", "kiran@cascade.test", "+1 555 0141", org(6), false},
		// Silverline Media
		{"Layla", "Jensen", "layla@silverline.test", "+1 555 0142", org(7), false},
		{"Miles", "Foster", "miles@silverline.test", "+1 555 0143", org(7), false},
		{"Nadia", "Becker", "nadia@silverline.test", "+1 555 0144", org(7), true},
		{"Owen", "Murray", "owen@silverline.test", "+1 555 0145", org(7), false},
		{"Pia", "Lindgren", "pia@silverline.test", "+1 555 0146", org(7), false},
		// Helio Systems
		{"Rohan", "Mehta", "rohan@helio.test", "+1 555 0147", org(8), false},
		{"Serena", "Costa", "serena@helio.test", "+1 555 0148", org(8), true},
		{"Tariq", "Aziz", "tariq@helio.test", "+1 555 0149", org(8), false},
		{"Ursula", "Braun", "ursula@helio.test", "+1 555 0150", org(8), false},
		{"Vince", "Russo", "vince@helio.test", "+1 555 0151", org(8), false},
		{"Wren", "Gallagher", "wren@helio.test", "+1 555 0152", org(8), false},
		// Polar Dynamics
		{"Xena", "Christou", "xena@polar.test", "+1 555 0153", org(9), false},
		{"Yusuf", "Demir", "yusuf@polar.test", "+1 555 0154", org(9), true},
		{"Zara", "Hussain", "zara@polar.test", "+1 555 0155", org(9), false},
		{"Arlo", "Bennett", "arlo@polar.test", "+1 555 0156", org(9), false},
		{"Bianca", "Fontana", "bianca@polar.test", "+1 555 0157", org(9), false},
		// Cedar Analytics
		{"Cyrus", "Hadid", "cyrus@cedar.test", "+1 555 0158", org(10), false},
		{"Dahlia", "Wu", "dahlia@cedar.test", "+1 555 0159", org(10), true},
		{"Elias", "Strand", "elias@cedar.test", "+1 555 0160", org(10), false},
		{"Flora", "Mancini", "flora@cedar.test", "+1 555 0161", org(10), false},
		{"Gustav", "Holm", "gustav@cedar.test", "+1 555 0162", org(10), false},
		// Orchid Capital
		{"Iris", "Takahashi", "iris@orchid.test", "+1 555 0163", org(11), false},
		{"Jasper", "Klein", "jasper@orchid.test", "+1 555 0164", org(11), true},
		{"Kaya", "Erdogan", "kaya@orchid.test", "+1 555 0165", org(11), false},
		{"Liam", "Barrett", "liam@orchid.test", "+1 555 0166", org(11), false},
		{"Maya", "Lund", "maya@orchid.test", "+1 555 0167", org(11), false},
		// Bluebird Software
		{"Noah", "Herrera", "noah@bluebird.test", "+1 555 0168", org(12), false},
		{"Olivia", "Janssen", "olivia@bluebird.test", "+1 555 0169", org(12), true},
		{"Pablo", "Ruiz", "pablo@bluebird.test", "+1 555 0170", org(12), false},
		{"Rosa", "Engström", "rosa@bluebird.test", "+1 555 0171", org(12), false},
		{"Simon", "Katz", "simon@bluebird.test", "+1 555 0172", org(12), false},
		// Mistral Consulting
		{"Thea", "Andersen", "thea@mistral.test", "+1 555 0173", org(13), false},
		{"Ulrich", "Bauer", "ulrich@mistral.test", "+1 555 0174", org(13), true},
		{"Vera", "Sousa", "vera@mistral.test", "+1 555 0175", org(13), false},
		{"Walter", "Ng", "walter@mistral.test", "+1 555 0176", org(13), false},
		{"Ximena", "Cruz", "ximena@mistral.test", "+1 555 0177", org(13), false},
		// Beacon Industries
		{"Yolanda", "Fischer", "yolanda@beacon.test", "+1 555 0178", org(14), false},
		{"Zayn", "Ali", "zayn@beacon.test", "+1 555 0179", org(14), true},
		{"Amara", "Okonkwo", "amara@beacon.test", "+1 555 0180", org(14), false},
		{"Boris", "Volkov", "boris@beacon.test", "+1 555 0181", org(14), false},
		{"Celine", "Morel", "celine@beacon.test", "+1 555 0182", org(14), false},
		// No organization
		{"Leo", "Park", "leo@example.test", "+1 555 0105", nil, false},
		{"Mila", "Novak", "mila@example.test", "+1 555 0183", nil, false},
		{"Nico", "Vega", "nico@example.test", "+1 555 0184", nil, true},
		{"Opal", "Stone", "opal@example.test", "+1 555 0185", nil, false},
		{"Perry", "Holmes", "perry@example.test", "+1 555 0186", nil, false},
		{"Riley", "Adams", "riley@example.test", "+1 555 0187", nil, false},
		{"Sage", "Turner", "sage@example.test", "+1 555 0188", nil, false},
		{"Talia", "Reed", "talia@example.test", "+1 555 0189", nil, false},
		{"Uri", "Gold", "uri@example.test", "+1 555 0190", nil, false},
		{"Val", "Stone", "val@example.test", "+1 555 0191", nil, false},
		{"Willa", "Grant", "willa@example.test", "+1 555 0192", nil, false},
		{"Xander", "Frost", "xander@example.test", "+1 555 0193", nil, false},
		{"Yael", "Shore", "yael@example.test", "+1 555 0194", nil, false},
		{"Zion", "Blake", "zion@example.test", "+1 555 0195", nil, false},
		{"Ada", "Sterling", "ada@example.test", "+1 555 0196", nil, false},
	}

	for _, seed := range contacts {
		_, err := database.CreateContact(db, database.Contact{
			OrganizationID: seed.orgID,
			FirstName:      seed.firstName,
			LastName:       seed.lastName,
			Email:          seed.email,
			Phone:          seed.phone,
			IsFavorite:     seed.favorite,
		})

		if err != nil {
			return err
		}
	}

	return nil
}

func seedNotes(db *sql.DB, now time.Time) error {
	notes := []struct {
		contactID int64
		userID    int64
		body      string
		ago       time.Duration
	}{
		{1, 1, "Followed up on onboarding checklist and confirmed next review date.", 2 * time.Hour},
		{1, 2, "Sent the revised contract for Alicia's review.", 18 * time.Hour},
		{1, 1, "Alicia confirmed she'll attend the quarterly sync.", 3 * 24 * time.Hour},
		{2, 1, "Requested revised pricing deck for the logistics pilot.", 26 * time.Hour},
		{2, 3, "Marcus needs the SLA appendix before Thursday.", 4 * 24 * time.Hour},
		{3, 1, "Juniper asked for a data residency addendum before procurement.", 49 * time.Hour},
		{3, 2, "Priya flagged a concern about the integration timeline.", 5 * 24 * time.Hour},
		{3, 4, "Sent compliance docs to Priya for internal review.", 8 * 24 * time.Hour},
		{4, 1, "Jules requested a demo of the analytics module.", 3 * time.Hour},
		{4, 3, "Confirmed the Summit Advisory renewal terms.", 6 * 24 * time.Hour},
		{5, 2, "Nora shared the Atlas Health security questionnaire.", 4 * time.Hour},
		{5, 1, "Completed the vendor assessment form for Atlas.", 7 * 24 * time.Hour},
		{6, 1, "Brandon asked about the enterprise tier pricing.", 5 * time.Hour},
		{7, 4, "Fatima confirmed the logistics integration timeline.", 6 * time.Hour},
		{7, 1, "Northstar team wants a dedicated support channel.", 10 * 24 * time.Hour},
		{8, 2, "George needs API docs for the warehouse module.", 8 * time.Hour},
		{9, 1, "Hannah requested a call to discuss renewal options.", 12 * time.Hour},
		{10, 3, "Isaac shared feedback on the beta release.", 15 * time.Hour},
		{12, 1, "Priya is excited about the new data pipeline feature.", 1 * time.Hour},
		{12, 2, "Set up a walkthrough for the Juniper team next week.", 20 * time.Hour},
		{13, 1, "Kai flagged a bug in the export module.", 24 * time.Hour},
		{14, 4, "Lena reviewed the compliance documentation.", 28 * time.Hour},
		{15, 1, "Mateo wants to expand their seat count.", 32 * time.Hour},
		{16, 2, "Nina confirmed attendance at the product roadmap session.", 36 * time.Hour},
		{17, 3, "Oscar asked for case studies from similar deployments.", 40 * time.Hour},
		{19, 1, "Jules wants to bring two more team members onboard.", 44 * time.Hour},
		{20, 1, "Quinn asked about the data migration process.", 2 * 24 * time.Hour},
		{21, 4, "Remi confirmed the Summit advisory board participation.", 52 * time.Hour},
		{22, 2, "Sana needs training materials for new hires.", 56 * time.Hour},
		{25, 1, "Nora wants to schedule a quarterly business review.", 60 * time.Hour},
		{26, 3, "Victor shared the updated compliance requirements.", 64 * time.Hour},
		{27, 1, "Wendy approved the contract amendment.", 68 * time.Hour},
		{28, 2, "Xavier flagged an issue with the billing module.", 72 * time.Hour},
		{31, 1, "Aiden requested a product roadmap preview.", 4 * 24 * time.Hour},
		{32, 4, "Bella confirmed the Meridian Finance partnership terms.", 100 * time.Hour},
		{33, 1, "Caleb asked about SSO integration options.", 104 * time.Hour},
		{34, 2, "Dina wants to pilot the analytics dashboard.", 108 * time.Hour},
		{37, 1, "Gael shared the Cascade Networks deployment plan.", 5 * 24 * time.Hour},
		{38, 3, "Hana confirmed the network audit schedule.", 128 * time.Hour},
		{39, 1, "Ivan needs access to the staging environment.", 132 * time.Hour},
		{42, 2, "Layla wants to explore the content management features.", 6 * 24 * time.Hour},
		{43, 4, "Miles asked for a bulk import template.", 152 * time.Hour},
		{44, 1, "Nadia confirmed the Silverline Media campaign launch.", 156 * time.Hour},
		{47, 1, "Rohan wants to discuss the Helio Systems upgrade path.", 7 * 24 * time.Hour},
		{48, 3, "Serena shared positive feedback from the pilot program.", 172 * time.Hour},
		{49, 2, "Tariq needs the API rate limits documentation.", 176 * time.Hour},
		{53, 1, "Xena asked about the disaster recovery SLA.", 8 * 24 * time.Hour},
		{54, 4, "Yusuf confirmed the Polar Dynamics contract renewal.", 196 * time.Hour},
		{55, 1, "Zara flagged a performance issue in the dashboard.", 200 * time.Hour},
		{59, 2, "Cyrus wants to schedule a Cedar Analytics demo.", 9 * 24 * time.Hour},
		{60, 1, "Dahlia shared the quarterly analytics report.", 220 * time.Hour},
		{61, 3, "Elias confirmed the data migration window.", 224 * time.Hour},
		{64, 1, "Iris asked about the Orchid Capital integration.", 10 * 24 * time.Hour},
		{65, 4, "Jasper confirmed the investment portfolio sync.", 244 * time.Hour},
		{66, 2, "Kaya wants custom reporting for their team.", 248 * time.Hour},
		{69, 1, "Noah shared the Bluebird Software release notes.", 11 * 24 * time.Hour},
		{70, 3, "Olivia confirmed the QA testing schedule.", 268 * time.Hour},
		{71, 1, "Pablo asked about the plugin architecture.", 272 * time.Hour},
		{74, 2, "Thea wants to discuss the Mistral Consulting scope.", 12 * 24 * time.Hour},
		{75, 4, "Ulrich confirmed the project kickoff date.", 292 * time.Hour},
		{76, 1, "Vera shared the client satisfaction survey results.", 296 * time.Hour},
		{79, 1, "Yolanda asked about the Beacon Industries roadmap.", 13 * 24 * time.Hour},
		{80, 3, "Zayn confirmed the manufacturing integration plan.", 316 * time.Hour},
		{81, 2, "Amara shared the onboarding documentation.", 320 * time.Hour},
		{84, 1, "Leo asked about community edition features.", 14 * 24 * time.Hour},
		{85, 4, "Mila wants to try the free tier before upgrading.", 340 * time.Hour},
		{86, 1, "Nico shared positive feedback about the UX redesign.", 344 * time.Hour},
		{1, 3, "Alicia's team completed the integration testing.", 15 * 24 * time.Hour},
		{6, 1, "Brandon confirmed the enterprise agreement.", 16 * 24 * time.Hour},
		{12, 4, "Priya wants to present at the next user conference.", 17 * 24 * time.Hour},
		{19, 2, "Jules shared the Summit Advisory case study draft.", 18 * 24 * time.Hour},
		{25, 1, "Nora completed the Atlas Health compliance review.", 19 * 24 * time.Hour},
		{31, 3, "Aiden asked about multi-region deployment.", 20 * 24 * time.Hour},
		{37, 1, "Gael confirmed the Cascade Networks go-live date.", 21 * 24 * time.Hour},
		{42, 4, "Layla shared the content strategy proposal.", 22 * 24 * time.Hour},
		{47, 2, "Rohan wants to explore the AI features.", 23 * 24 * time.Hour},
		{53, 1, "Xena confirmed the Polar Dynamics disaster recovery test.", 24 * 24 * time.Hour},
		{59, 3, "Cyrus completed the Cedar Analytics onboarding.", 25 * 24 * time.Hour},
		{64, 1, "Iris shared the Orchid Capital quarterly review.", 26 * 24 * time.Hour},
		{69, 4, "Noah confirmed the Bluebird Software sprint plan.", 27 * 24 * time.Hour},
		{74, 2, "Thea wants custom training for Mistral Consulting.", 28 * 24 * time.Hour},
		{79, 1, "Yolanda confirmed the Beacon Industries pilot results.", 29 * 24 * time.Hour},
		{2, 1, "Marcus asked for an updated shipping integration guide.", 30 * time.Hour},
		{4, 2, "Jules wants to schedule a strategy session.", 50 * time.Hour},
		{7, 3, "Fatima shared the warehouse optimization report.", 70 * time.Hour},
		{9, 1, "Hannah confirmed the annual review meeting.", 90 * time.Hour},
		{11, 4, "Jana asked about the bulk processing module.", 110 * time.Hour},
		{14, 1, "Lena completed the compliance audit checklist.", 130 * time.Hour},
		{16, 2, "Nina wants to demo the platform to her leadership.", 150 * time.Hour},
		{18, 3, "Paloma asked about the API webhook capabilities.", 170 * time.Hour},
		{21, 1, "Remi confirmed the advisory board meeting agenda.", 190 * time.Hour},
		{23, 4, "Tomás shared the project timeline update.", 210 * time.Hour},
		{26, 1, "Victor asked about HIPAA compliance features.", 230 * time.Hour},
		{29, 2, "Yara wants to explore the mobile app integration.", 250 * time.Hour},
		{32, 1, "Bella confirmed the quarterly financials review.", 270 * time.Hour},
		{35, 3, "Ethan shared the market analysis report.", 290 * time.Hour},
		{38, 1, "Hana confirmed the network security assessment.", 310 * time.Hour},
		{41, 4, "Kiran asked about the custom dashboard builder.", 330 * time.Hour},
		{44, 1, "Nadia wants to present the media campaign results.", 350 * time.Hour},
		{48, 2, "Serena shared the pilot program metrics.", 370 * time.Hour},
		{50, 1, "Ursula asked about the support tier options.", 390 * time.Hour},
		{54, 3, "Yusuf confirmed the annual contract renewal.", 410 * time.Hour},
		{57, 1, "Arlo asked about edge computing integration.", 430 * time.Hour},
		{60, 4, "Dahlia shared the data governance framework.", 450 * time.Hour},
		{63, 1, "Gustav confirmed the analytics pipeline upgrade.", 470 * time.Hour},
		{66, 2, "Kaya wants to schedule an executive briefing.", 490 * time.Hour},
		{70, 1, "Olivia confirmed the QA automation rollout.", 510 * time.Hour},
		{73, 3, "Rosa asked about the localization features.", 530 * time.Hour},
		{76, 1, "Vera shared the customer success metrics.", 550 * time.Hour},
		{80, 4, "Zayn confirmed the manufacturing module update.", 570 * time.Hour},
		{83, 1, "Celine asked about the French language support.", 590 * time.Hour},
	}

	for _, note := range notes {
		if _, err := database.CreateNoteAt(db, note.contactID, note.userID, note.body, now.Add(-note.ago)); err != nil {
			return err
		}
	}

	return nil
}
