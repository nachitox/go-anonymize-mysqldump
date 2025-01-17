package main

import (
	"bytes"
	"github.com/omnea/faker"
	"github.com/stretchr/testify/assert"
	"testing"
)

var (
	jsonConfig         Config
	dropAndCreateTable = "DROP TABLE IF EXISTS `wp_options`;\n" +
		"/*!40101 SET @saved_cs_client     = @@character_set_client */;\n" +
		"/*!40101 SET character_set_client = utf8 */;\n" +
		"CREATE TABLE `wp_options` (\n" +
		"`option_id` bigint(20) unsigned NOT NULL AUTO_INCREMENT,\n" +
		"`option_name` varchar(191) NOT NULL DEFAULT '',\n" +
		"`option_value` longtext NOT NULL,\n" +
		"`autoload` varchar(20) NOT NULL DEFAULT 'yes',\n" +
		"PRIMARY KEY (`option_id`),\n" +
		"UNIQUE KEY `option_name` (`option_name`)\n" +
		") ENGINE=InnoDB AUTO_INCREMENT=123 DEFAULT CHARSET=utf8mb4;\n" +
		"/*!40101 SET character_set_client = @saved_cs_client */;"

	// Don't forget to escape \ because it'll translate to a newline and not pass
	// the comparison test
	multilineQuery = `INSERT INTO wp_usermeta VALUES
	(1,1,'first_name','John'),(2,1,'last_name','Doe'),
	(3,1,'foobar','bazquz'),
	(4,1,'nickname','Jim'),
	(5,1,'description','Lorum ipsum.');`
	multilineQueryRecompiled = "insert into wp_usermeta values (1, 1, 'first_name', 'Libero ea est.'), (2, 1, 'last_name', 'Sequi cum expedita.'), (3, 1, 'foobar', 'Minus omnis est.'), (4, 1, 'nickname', 'Aliquam nam dolores.'), (5, 1, 'description', 'Asperiores est nesciunt.');\n"
	commentsQuery            = "INSERT INTO `wp_comments` VALUES (1,1,'A WordPress Commenter','wapuu@wordpress.example','https://wordpress.org/','','2019-06-12 00:59:19','2019-06-12 00:59:19','Hi, this is a comment.\\nTo get started with moderating, editing, and deleting comments, please visit the Comments screen in the dashboard.\\nCommenter avatars come from <a href=\\\"https://gravatar.com\\\">Gravatar</a>.',0,'1','','',0,0);\n"
	commentsQueryRecompiled  = "insert into wp_comments values (1, 1, 'elsa', 'nisa.ksters@example.net', 'http://streck.ch/lion', '', '2019-06-12 00:59:19', '2019-06-12 00:59:19', 'Hi, this is a comment.\\nTo get started with moderating, editing, and deleting comments, please visit the Comments screen in the dashboard.\\nCommenter avatars come from <a href=\\\"https://gravatar.com\\\">Gravatar</a>.', 0, '1', '', '', 0, 0);\n"
	usersQuery               = "INSERT INTO `wp_users` VALUES (1,'username','user_pass','username','hosting@humanmade.com','','2019-06-12 00:59:19','dc316129015f6782eac32463c83bec40fa0b6be6',0,'username'),(2,'username','user_pass','username','hosting@humanmade.com','http://notreal.com/username','2019-06-12 00:59:19','9195de087ee7c5404e8102c2ebd3ed21b3d253f7',0,'username');\n"
	usersQueryRecompiled     = "insert into wp_users values (1, 'josie_schoberg', 'NjaK5HeMAMuv', 'mona', 'teo.molzan@example.net', '', '2019-06-12 00:59:19', 'UAK-0001', 0, 'Dilara Wilky'), (2, 'sinja', 'J3JRQ4XoIxXX6A', 'luka.rapp', 'abby.wachenbrunner@example.net', 'http://porth.com/greta', '2019-06-12 00:59:19', 'UAK-0001', 0, 'Sascha von Lindner');\n"
	userMetaQuery            = "INSERT INTO `wp_usermeta` VALUES (1,1,'first_name','John'),(2,1,'last_name','Doe'),(3,1,'foobar','bazquz'),(4,1,'nickname','Jim'),(5,1,'description','Lorum ipsum.'),(6,2,'first_name','Janet'),(7,2,'last_name','Doe'),(8,2,'foobar','bazquz'),(9,2,'nickname','Jane'),(10,2,'description','Lorum ipsum.');\n"
	userMetaQueryRecompiled  = "insert into wp_usermeta values (1, 1, 'first_name', 'Dolorum nostrum alias.'), (2, 1, 'last_name', 'Qui voluptatum est.'), (3, 1, 'foobar', 'Eveniet repellat in.'), (4, 1, 'nickname', 'Eligendi quia ex.'), (5, 1, 'description', 'Consequuntur dolores facilis.'), (6, 2, 'first_name', 'Facilis ut unde.'), (7, 2, 'last_name', 'Quisquam enim consequatur.'), (8, 2, 'foobar', 'Unde velit reiciendis.'), (9, 2, 'nickname', 'Eaque est reiciendis.'), (10, 2, 'description', 'Voluptas eum consequatur.');\n"
)

func init() {
	faker.Seed(432)
	jsonConfig = readConfigFile("./config.example.json")
}

func BenchmarkProcessLine(b *testing.B) {
	for i := 0; i < b.N; i++ {
		processLine(usersQuery, jsonConfig)
		processLine(userMetaQuery, jsonConfig)
		processLine(commentsQuery, jsonConfig)
	}
}

func TestSetupAndProcessInput(t *testing.T) {

	var tests = []struct {
		testName string
		query    string
		wants    string
	}{
		{
			testName: "users query",
			query:    usersQuery,
			wants:    usersQueryRecompiled,
		},
		{
			testName: "usermeta query",
			query:    userMetaQuery,
			wants:    userMetaQueryRecompiled,
		},
		{
			testName: "comments query",
			query:    commentsQuery,
			wants:    commentsQueryRecompiled,
		},
		{
			testName: "multiline query",
			query:    multilineQuery,
			wants:    multilineQueryRecompiled,
		},
		{
			testName: "table creation",
			query:    dropAndCreateTable,
			wants:    dropAndCreateTable + "\n",
		},
	}

	for _, test := range tests {
		t.Run(test.testName, func(t *testing.T) {

			input := bytes.NewBufferString(test.query)
			lines := setupAndProcessInput(jsonConfig, input)

			var result string
			for line := range lines {
				result += <-line
			}

			if result != test.wants {
				t.Error("\nExpected:\n", test.wants, "\nActual:\n", result)
			}
		})
	}
}

func TestUniqueMap(t *testing.T) {
	setCurrentTable("t")
	setCurrentField("field")

	assert.False(t, checkMapExists("field", "value"), "Value should not exist")

	setMapValue("field", "value")
	assert.True(t, checkMapExists("field", "value"), "Value should exist now")
}
