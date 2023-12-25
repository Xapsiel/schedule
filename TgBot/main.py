import telebot, datetime,json
import requests, fnmatch

from telebot import types

bot = telebot.TeleBot('6479027249:AAGwRywi8194mHVIsAZ3VQCvmKvDIx3XD6o')
is_student = "unknown"
name = "СмирноваС.И."
group = "233-1"
def delete_message(message):
    chat_id = message.chat.id
    message_id = message.message_id
    bot.delete_message(chat_id, message_id)
def schedule_for_teacher(year,month,day, name, message):
    response = requests.get("http://localhost:80//schedule/teacher/pairInfo/"+str(year)+"/"+str(month)+"/"+str(day)+"/"+name)
    if response.text =='':
        bot.send_message(message.chat.id, text="У вас сегодня нет пар")
        return
    c = json.loads(response.text)
    if c == "null": 
        bot.send_message(chat_id=message.chat.id, text="Сегодня пар нет!\n^_^\n")
        return
    lessons = c["Lessons"]
    for k,v in lessons.items():
        if v[0][3]=='лк': v[0][3]='Лекция'
        else: v[0][3]='Практическое занятие'
        if v[0][4]=="":v[0][4]="Комментария нет"
        result = f"\t*Пара {k}:*\n*Название предмета:*{v[0][0]}\n*Аудитория:*{v[0][2]}\n*Тип занятия:*{v[0][3]}\n*Ваш комментарий:*{v[0][4]}\n"
        result += "*Занятие у этих(ой) групп(ы):*\n"
        if len(v)!=0:
            for info in v:
                result+=f"\t{info[1]}\n"
        bot.send_message(message.chat.id, text=result,parse_mode="Markdown")
def schedule_for_student(year,month,day, group, message):
    lesson = {
        0:"Название предмета: ",
        1:"Имя преподавателя: ",
        2:"Адрес: ",
        3:"Тип пары: ",
        4:"Комментарий преподавателя: ",
        5:"Время начала пары: "
    }
    if int(day) < 10:
        day = "0"+str(day)
    if month<10:
        month = "0"+str(month)
    response = requests.get("http://localhost:80//schedule/student/pairInfo/"+str(year)+"/"+str(month)+"/"+str(day)+"/"+group)
    if response.text =='':
        bot.send_message(message.chat.id, text="Нет пар")
        return
    c = json.loads(response.text)
    if c == "null": 
        bot.send_message(chat_id=message.chat.id, text="Нет пар\n^_^\n")
        return
    information = [i[1] for i in c["Lessons"].items()]
    for i in range(0,len(information)):
        text = information[i]
        if len(text[0])== 0:
            continue
        result = ""
        if text[3] == "лк": 
            text[3]="Лекция"
        else: 
            text[3]="Практическое занятие"
        if text[4]=="": text[4]=="Комментария нет"
        for col in range(len(lesson)):
            result += "*"+lesson[col]+"*"+text[col]+"\n"
        number = i+1
        bot.send_message(chat_id=message.chat.id, text=str(number)+" пара:\n"+result,parse_mode="Markdown")
@bot.message_handler(commands=['start'])
def start(message) :
    markup = types.ReplyKeyboardMarkup(resize_keyboard=True)
    i_am_teacher = types.KeyboardButton("Я преподаватель")
    i_am_student = types.KeyboardButton("Я студент")
    markup.add(i_am_student,i_am_teacher)
    bot.send_message(message.chat.id, text="Кто вы?", reply_markup=markup)

@bot.message_handler(content_types="text")
def message_reply(message):
    
    if  message.text == "Вернуться в меню":
        if globals()["is_student"]:
            message.text = globals()["group"]
        else: message.text = globals()["name"]
    if message.text == "Пройти регистрацию":
        message.text = "Регистрация"
    if message.text == 'Регистрация' :
        markup = types.ReplyKeyboardMarkup(resize_keyboard = True)
        i_am_teacher = types.KeyboardButton("Я преподаватель")
        i_am_student = types.KeyboardButton("Я студент")
        markup.add(i_am_student,i_am_teacher)
        bot.send_message(message.chat.id, text="Кто вы?", reply_markup=markup)
    
    elif message.text=="Я студент" :
        globals()["is_student"] = True
        bot.send_message(message.chat.id, text='Введите вашу группу в формате "группа-подгрупп(233-1)"')
    elif fnmatch.fnmatch(message.text.lower(), "23[1-3]-[1-2]") and globals()["is_student"]:
        globals()["group"] = message.text
        markup = types.ReplyKeyboardMarkup(resize_keyboard=True)
        today = types.KeyboardButton("Расписание на сегодня")
        tomorrow = types.KeyboardButton("Расписание на завтра")
        to_weekday = types.KeyboardButton("Расписание на <???>")
        next_pair = types.KeyboardButton("Следующая пара")
        to_menu = types.KeyboardButton("Пройти регистрацию")
        markup.add(today,tomorrow,to_weekday,next_pair,to_menu)
        bot.send_message(message.chat.id, text='Что от меня требуется?', reply_markup=markup)

    elif globals()["is_student"] and (message.text == "Расписание на сегодня" or message.text == "Расписание на завтра"):
        date = datetime.datetime.now()
        if message.text == "Расписание на завтра":
            date+=datetime.timedelta(days=1)
        year = date.year
        month = date.month
        day = date.day
        if day < 10:
            day = "0"+str(day)
        if month<10:
            month = "0"+str(month)
        schedule_for_student(year=year,month=month,day=day,group=group,message=message)
    elif message.text == "Расписание на <???>":
        markup = types.ReplyKeyboardMarkup(resize_keyboard=True)
        monday = types.KeyboardButton("Понедельник")
        tuesday = types.KeyboardButton("Вторник")
        wednesday = types.KeyboardButton("Среда")
        thursday = types.KeyboardButton("Четверг")
        friday = types.KeyboardButton("Пятница")
        to_menu = types.KeyboardButton("Вернуться в меню")

        markup.add(monday,tuesday,wednesday,thursday,friday,to_menu)
        bot.send_message(message.chat.id, text="Расписание на какой день недели вам нужно?", reply_markup=markup)
    elif globals()["is_student"] and message.text in ["Понедельник","Вторник","Среда","Четверг","Пятница"]:
        weekday = message.text.lower()
        from_week_to_num = {
            "понедельник":0,
            "вторник":1,
            "среда":2,
            "четверг":3,
            "пятница":4,
            
        }
        weekday = from_week_to_num[weekday]

        now_day = datetime.datetime.now() 
        while now_day.weekday()!=weekday:
            now_day += datetime.timedelta(days=1)
        year = now_day.year
        month = now_day.month
        day = now_day.day
        schedule_for_student(year=year,month=month,day=day, group=group,message=message)
    elif globals()["is_student"] and message.text == "Следующая пара":
        response = requests.get("http://localhost:80///schedule/student/next/"+group)
        information = response.text
        if information == 'null':
            bot.send_message(message.chat.id,text="Следующая пара-окно")
        else:
            lesson = {
                0:"Название предмета: ",
                1:"Имя преподавателя: ",
                2:"Адрес: ",
                3:"Тип пары: ",
                4:"Комментарий преподавателя: ",
                5:"Время начала пары: "
            }   
             
            text = json.loads(information )
            if len(text[0])== 0:
                bot.send_message(message.chat.id,text="Следующая пара-окно")
                return
            
            result = ""
            if text[3] == "лк": 
                text[3]="Лекция"
            else: 
                text[3]="Практическое занятие"
            if text[4]=="": text[4]="Комментария нет"
            for col in range(len(lesson)):
                result += "*"+lesson[col]+"*"+text[col]+"\n"
            
            bot.send_message(chat_id=message.chat.id, text= result,parse_mode="Markdown")

    elif message.text=="Я преподаватель":
        globals()["is_student"] = False
        bot.send_message(message.chat.id, text="Введите ваши инициалы в формате <Иванов И.П>(Иванов Иван Петрович)")
    elif message.text.count(".")==2 and not globals()["is_student"]:
        globals()["name"] = ''.join(message.text.split())
        markup = types.ReplyKeyboardMarkup(resize_keyboard=True)
        today = types.KeyboardButton("Расписание на сегодня")
        tomorrow = types.KeyboardButton("Расписание на завтра")
        to_weekday = types.KeyboardButton("Расписание на <???>")
        next_pair = types.KeyboardButton("Следующая пара")
        to_menu = types.KeyboardButton("Пройти регистрацию")
        markup.add(today,tomorrow,to_weekday,next_pair,to_menu)
        bot.send_message(message.chat.id, text="Что от меня требуется?", reply_markup=markup)
    elif (message.text == "Расписание на сегодня" or message.text == "Расписание на завтра")  and not globals()["is_student"]:
        date = datetime.datetime.now()
        if message.text == "Расписание на завтра":
            date+=datetime.timedelta(days=1)
        year = date.year
        month = date.month
        day = date.day
        if day < 10:
            day = "0"+str(day)
        if month<10:
            month = "0"+str(month)
        schedule_for_teacher(year=year,month=month,day=day,name=name,message=message)
        
     
        
    elif not globals()["is_student"] and message.text in ["Понедельник","Вторник","Среда","Четверг","Пятница"]:
        weekday = message.text.lower()
        from_week_to_num = {
            "понедельник":0,
            "вторник":1,
            "среда":2,
            "четверг":3,
            "пятница":4,
            
        }
        weekday = from_week_to_num[weekday]

        now_day = datetime.datetime.now() 
        while now_day.weekday()!=weekday:
            now_day += datetime.timedelta(days=1)
        year = now_day.year
        month = now_day.month
        day = now_day.day
        schedule_for_teacher(year=year,month=month,day=day, name=name,message=message)
    
    elif message.text == "Следующая пара" and not globals()["is_student"]:
        response = requests.get("http://localhost:80//schedule/teacher/next/"+name)
        information =  json.loads(response.text)
        if information is None: 
            bot.send_message(message.chat.id, text="Следующая пара-окно")
            return
        if information[0][3]=="лк":information[0][3]="Лекция"
        else: information[0][3]="Практическое занятие"
        result = f"*Название предмета: *{information[0][0]}\n*Адрес: *{information[0][2]}\n*Тип занятия: *{information[0][3]}\n*Ваш комментарий: *{information[0][4]}\n*Занятие у этих(ой) групп(ы):*\n"
        for elem in information:
            result += f"{elem[1]}"
        bot.send_message(message.chat.id,result,parse_mode="Markdown")
  
bot.polling(none_stop=True)