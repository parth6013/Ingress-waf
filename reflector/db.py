import sqlite3


def connect_db():
    db = sqlite3.connect("database.db")
    db.cursor().execute(
        "CREATE TABLE IF NOT EXISTS User_Details"
        "(id INTEGER PRIMARY KEY, "
        "username TEXT,"
        "salary TEXT,"
        "grade TEXT,"
        "DOB TEXT,"
        "phone_number TEXT)"
    )
    db.commit()
    return db


def add_comment(username, salary, grade, DOB, phone_number):
    db = connect_db()
    db.cursor().execute(
        "INSERT INTO User_Details (username, salary, grade, DOB, phone_number) "
        "VALUES (?, ?, ?, ?, ?)",
        (username, salary, grade, DOB, phone_number),
    )
    db.commit()


def get_comments(search_query):

  
    db = connect_db()
    results = []
    try:
        get_all_query = (
            "SELECT username, salary, grade, DOB, phone_number FROM User_Details WHERE id="
            + str(search_query)
        )
        for username, salary, grade, DOB, phone_number in (
            db.cursor().execute(get_all_query).fetchall()
        ):
            # Formatted output
            formatted_result = (
                f"NAME: {username} "
                f"Phone_number: {phone_number} "
                f"Salary: {salary} "
                f"Grade: {grade} "
                f"DOB: {DOB}"
            )
            results.append(formatted_result)
    except sqlite3.OperationalError:
        pass

    return results





  # add_comment("john_doe", "55000", "B", "1990-01-15", "123-456-7890")
    # add_comment("jane_smith", "60000", "A", "1985-06-22", "234-567-8901")
    # add_comment("alex_jones", "52000", "C", "1992-03-11", "345-678-9012")
    # add_comment("sara_brown", "70000", "B+", "1988-09-05", "456-789-0123")
    # add_comment("mike_davis", "48000", "B-", "1995-12-30", "567-890-1234")
    # print(get_comments("1 OR 1=1"))