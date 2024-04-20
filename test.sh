#!/bin/bash

# Array of possible characters for random name generation
characters=(a b c d e f g h i j k l m n o p q r s t u v w x y z)

# Function to generate a random string of specified length
generate_random_string() {
    local length=$1
    local random_string=""
    for (( i=0; i<$length; i++ )); do
        random_index=$((RANDOM % ${#characters[@]}))
        random_string+=${characters[random_index]}
    done
    echo "$random_string"
}

# Loop to send random requests
for (( i=0; i<1000; i++ )); do
    # Generate random lengths for "name"
    lengths=(3 6 60)
    random_length=${lengths[$RANDOM % ${#lengths[@]}]}
    name=$(generate_random_string $random_length)

    # Generate random user-agent
    user_agents=("Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/58.0.3029.110 Safari/537.3" \
                "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/72.0.3626.121 Safari/537.36" \
                "Mozilla/5.0 (Windows NT 10.0; Win64; x64; rv:85.0) Gecko/20100101 Firefox/85.0")
    random_user_agent=${user_agents[$RANDOM % ${#user_agents[@]}]}

    # Generate random delay between 100 ms and 2 seconds
    random_delay=$(( ( RANDOM % 1900 ) + 100 ))

    # Send request using curl with delay
    curl -A "$random_user_agent" -X GET "http://localhost:80/?name=$name" >/dev/null 2>&1 &

    # Sleep for random delay
    sleep $(bc <<< "scale=2; $random_delay / 1000")
done

echo "Finish"