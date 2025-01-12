# Game loop

 - Player joins with a certain balance
 - Player bets a certain amount
 - Player receives a hand if they have bet something
 - If everyone has bet or the time limit is reached
	 - All players receive a card turn for turn (including the dealer)
	 - Until all players have 2 cards (including the dealer (one closed card))
	 - All hands are checked for blackjack
	 	- If someone has blackjack, winnings are paid (the payout is 3-to-2, this means if the bet is 10 the payout is 15 + the orignal 10)
	 - For all players:
		- Player gets to choose if they want to hit
			- If player chooses to hit
				- Player gets an extra card
				- If new total of the cards is bigger than 21
					- Player loses
			- If player chooses not to hit
				- Next player turn
	 - After everyone has had their turn
	 	- Dealer hits until total is greater than 16 (so 17 or higher)
			- If dealer gets a total higher than 21 (so 22 or higher)
				- All players that are is still in the game win and receive double their bet
	 - If the player has a total lower than the dealer
	 	- Player loses
	 - If the player has a total equal to the dealer
	 	- Player gets their bet back
	 - If the player has a total higher than the dealer
	 	- Player gets double their bet

 # Testing

 - Game State logic
 - Player hit/stand
 - Ending of game
 - Dealer turns
 - Starting new games
 - Players leaving

# TO-DO

 - Blackjack detection after dealing cards (Player & Dealer)
 	- Including payouts
 - Splitting hands
 - Database integration
 	- Database interface for playing along an ongoing game in the DB
	- SQLite implementation
 - REST API
 	- WebSocket for messages to the server
	- SSE for updates in game state
