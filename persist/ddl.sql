drop table if exists question;
drop table if exists answer;
drop table if exists question_answer;
drop table if exists session;

create table if not exists question (
	id integer primary key,
	body text not null,
	number integer not null,
	correct_index integer not null,
	info text not null
);

create table if not exists answer (
	id integer primary key,
	body text not null,
	answer_index integer not null,
	question_id integer not null,
	foreign key(question_id) references question(id)
);

create table if not exists question_answer (
	id integer primary key,
	session_id text not null,
	answer_index integer not null,
	question_id integer not null,
	foreign key(question_id) references question(id)
);

insert into question (body, number, correct_index, info) values ("In parts per million (ppm), what is the concentration of carbon dioxide in the atmosphere today?", 0, 3, "Carbon dioxide concentration in the atmosphere was measured at 280 ppm before industrialization. This is an increase in over 120 ppm, or 44.9%, in about 200 years.");
insert into question (body, number, correct_index, info) values ("How much carbon dioxide does one gallon of gasoline emit into the atmosphere?", 1, 2, "The United States goes went through 143.37 billion gallons of gasoline in 2016. That’s nearly 3 trillion pounds of CO2 emissions in the US, in 2016, alone.");
insert into question (body, number, correct_index, info) values ("Glacier National Park used to be home to an estimated 150 glaciers. How many glaciers are in the park today?", 2, 0, "By 1968, Glacier National Park was home to only 50 glaciers. It is expected that it will lose all of them by 2030 due to warming waters.");
insert into question (body, number, correct_index, info) values ("The United States is home to roughly 326,625,791 people, or 4.4% of the world’s population. What share of global CO2 emissions does it produce?", 3, 1, "China is the only country that emits more CO2 into the atmosphere than the United States. China’s share of global emissions is a whopping 23.43%.");
insert into question (body, number, correct_index, info) values ("How many acres of rainforest are destroyed every day?", 4, 3, "Given that a single tree can convert as much as 48 pounds of carbon dioxide to clean air per year, it is no wonder CO2 emissions are quickly becoming a problem.");
insert into question (body, number, correct_index, info) values ("Since the late 19th century, how much has the average surface temperature of the Earth increased?", 5, 1, "Most of the warming has occurred within the past 35 years. 16 of the 17 hottest recorded years have occurred since 2001, with 2016 being the hottest year ever recorded.");
insert into question (body, number, correct_index, info) values ("How much of the Earth’s coral reefs has been impacted by coral bleaching?", 6, 2, "Coral bleaching occurs when corals are exposed to stressful conditions, like high temperature. Corals expel the symbiotic algae living in their tissues, causing corals to turn white or pale. Without the algae, the coral loses its major source of food and is more susceptible to disease. 4,630 square miles of reef have been killed off due to bleaching.");
insert into question (body, number, correct_index, info) values ("How much are sea levels projected to rise by 2100?", 7, 3, "Sea levels globally have risen 8 inches in the last century, and will continue to rise 3.4 millimeters every year.");
insert into question (body, number, correct_index, info) values ("Cumulatively in 2016 and 2017, how much impact did FGCU’s participation in RecycleMania have on greenhouse gas reduction?", 8, 2, "This is equal to removing 38 cars off the road, or eliminating the energy consumption of 17 households.");
insert into question (body, number, correct_index, info) values ("How many pounds of recycling did FGCU contribute to RecycleMania in 2017?", 9, 0, "That’s almost 14,000 pounds of recycling every week during the competition, or 7.47 pounds contributed per student during the competition.");

insert into answer (body, answer_index, question_id) values ("329.1 ppm", 0, 1);
insert into answer (body, answer_index, question_id) values ("352.6 ppm", 1, 1);
insert into answer (body, answer_index, question_id) values ("384.5 ppm", 2, 1);
insert into answer (body, answer_index, question_id) values ("405.6 ppm", 3, 1);

insert into answer (body, answer_index, question_id) values ("10.54 pounds", 0, 2);
insert into answer (body, answer_index, question_id) values ("15.63 pounds", 1, 2);
insert into answer (body, answer_index, question_id) values ("19.64 pounds", 2, 2);
insert into answer (body, answer_index, question_id) values ("21.98 pounds", 3, 2);

insert into answer (body, answer_index, question_id) values ("25", 0, 3);
insert into answer (body, answer_index, question_id) values ("40", 1, 3);
insert into answer (body, answer_index, question_id) values ("55", 2, 3);
insert into answer (body, answer_index, question_id) values ("70", 3, 3);
	
insert into answer (body, answer_index, question_id) values ("9.35%", 0, 4);
insert into answer (body, answer_index, question_id) values ("14.69%", 1, 4);
insert into answer (body, answer_index, question_id) values ("17.32%", 2, 4);
insert into answer (body, answer_index, question_id) values ("20.80%", 3, 4);
	
insert into answer (body, answer_index, question_id) values ("26,000", 0, 5);
insert into answer (body, answer_index, question_id) values ("38,000", 1, 5);
insert into answer (body, answer_index, question_id) values ("50,000", 2, 5);
insert into answer (body, answer_index, question_id) values ("67,000", 3, 5);
	
insert into answer (body, answer_index, question_id) values ("1° F", 0, 6);
insert into answer (body, answer_index, question_id) values ("2° F", 1, 6);
insert into answer (body, answer_index, question_id) values ("3° F", 2, 6);
insert into answer (body, answer_index, question_id) values ("4° F", 3, 6);
	
insert into answer (body, answer_index, question_id) values ("25%", 0, 7);
insert into answer (body, answer_index, question_id) values ("30%", 1, 7);
insert into answer (body, answer_index, question_id) values ("40%", 2, 7);
insert into answer (body, answer_index, question_id) values ("50%", 3, 7);
	
insert into answer (body, answer_index, question_id) values ("0.11 to 1.1 feet", 0, 8);
insert into answer (body, answer_index, question_id) values ("0.33 to 3.3 feet", 1, 8);
insert into answer (body, answer_index, question_id) values ("0.44 to 4.4 feet", 2, 8);
insert into answer (body, answer_index, question_id) values ("0.66 to 6.6 feet", 3, 8);

insert into answer (body, answer_index, question_id) values ("94 metric tons of CO2 equivalent", 0, 9);
insert into answer (body, answer_index, question_id) values ("127 metric tons of CO2 equivalent", 1, 9);
insert into answer (body, answer_index, question_id) values ("193 metric tons of CO2 equivalent", 2, 9);
insert into answer (body, answer_index, question_id) values ("218 metric tons of CO2 equivalent", 3, 9);

insert into answer (body, answer_index, question_id) values ("110,925", 0, 10);
insert into answer (body, answer_index, question_id) values ("126,325", 1, 10);
insert into answer (body, answer_index, question_id) values ("145,800", 2, 10);
insert into answer (body, answer_index, question_id) values ("163,245", 3, 10);


