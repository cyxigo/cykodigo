package data

import "github.com/bwmarrin/discordgo"

const (
	CmdHelp        = "help"
	CmdMe          = "me"
	CmdMeow        = "meow"
	CmdMeowat      = "meowat"
	CmdBark        = "bark"
	CmdBarkat      = "barkat"
	CmdFlip        = "flip"
	CmdExplode     = "explode"
	CmdSpin        = "spin"
	CmdDance       = "dance"
	CmdCat         = "cat"
	CmdCart        = "cart"
	CmdHate        = "hate"
	CmdLove        = "love"
	CmdAssault     = "assault"
	CmdBalance     = "balance"
	CmdBalanceall  = "balanceall"
	CmdShop        = "shop"
	CmdInventory   = "inventory"
	CmdLeaderboard = "leaderboard"
	CmdRoulette    = "roulette"
	CmdWork        = "work"
	CmdTransfer    = "transfer"
	CmdSteal       = "steal"
	CmdBuy         = "buy"
	CmdSell        = "sell"
	CmdGive        = "give"
	CmdEat         = "eat"
	CmdCooldown    = "cooldown"
	CmdHigh        = "high"

	// msg means that this is message content command
	// same for handlers: handleMsg
	CmdMsgMeow         = "meow"
	CmdMsgCrazy        = "crazy"
	CmdMsgExplodeBalls = "cykodigo explode balls"
	CmdMsgGlamptastic  = "glamptastic"
)

// slash commands, looooooooooong list of them
var Cmds = []*discordgo.ApplicationCommand{
	{
		Name:        CmdHelp,
		Description: "This... will actually help, well, maybe",
	},
	{
		Name:        CmdMe,
		Description: "Send picture of me! Nah, not you of course",
	},
	{
		Name:        CmdMeow,
		Description: "He will meow",
	},
	{
		Name:        CmdMeowat,
		Description: "Meow at someone",
		Options: []*discordgo.ApplicationCommandOption{
			{
				Type:        discordgo.ApplicationCommandOptionUser,
				Name:        "user",
				Description: "Person to meow at",
				Required:    true,
			},
		},
	},
	{
		Name:        CmdBark,
		Description: "He will... or won't bark. I don't know he's a cat",
	},
	{
		Name:        CmdBarkat,
		Description: "Bark at someone",
		Options: []*discordgo.ApplicationCommandOption{
			{
				Type:        discordgo.ApplicationCommandOptionUser,
				Name:        "user",
				Description: "Person to bark at",
				Required:    true,
			},
		},
	},
	{
		Name:        CmdFlip,
		Description: "He will do a flip",
	},
	{
		Name:        CmdExplode,
		Description: "He will explo- WHAT??? DON'T DO THAT!!!",
	},
	{
		Name:        CmdSpin,
		Description: "He will spin! Wooooooo",
	},
	{
		Name:        CmdDance,
		Description: "Let's dance",
	},
	{
		Name:        CmdCat,
		Description: "Cat!",
	},
	{
		Name:        CmdCart,
		Description: "Cart!",
	},
	{
		Name:        CmdHate,
		Description: "Hate someone",
		Options: []*discordgo.ApplicationCommandOption{
			{
				Type:        discordgo.ApplicationCommandOptionUser,
				Name:        "user",
				Description: "Person to hate",
				Required:    true,
			},
		},
	},
	{
		Name:        CmdLove,
		Description: "Love someone",
		Options: []*discordgo.ApplicationCommandOption{
			{
				Type:        discordgo.ApplicationCommandOptionUser,
				Name:        "user",
				Description: "Person to love",
				Required:    true,
			},
		},
	},
	{
		Name:        CmdAssault,
		Description: "Try to assault someone... shhh...",
		Options: []*discordgo.ApplicationCommandOption{
			{
				Type:        discordgo.ApplicationCommandOptionUser,
				Name:        "user",
				Description: "Person to try to assault",
				Required:    true,
			},
			{
				Type:        discordgo.ApplicationCommandOptionString,
				Name:        "item",
				Description: "Item to assault someone with",
				Required:    true,
			},
		},
	},
	{
		Name:        CmdBalance,
		Description: "Check how much money you've got from hard work",
		Options: []*discordgo.ApplicationCommandOption{
			{
				Type:        discordgo.ApplicationCommandOptionUser,
				Name:        "user",
				Description: "Person whose balance to show (optional)",
				Required:    false,
			},
		},
	},
	{
		Name:        CmdBalanceall,
		Description: "Check how much money each member of the server has",
	},
	{
		Name:        CmdShop,
		Description: "View available items",
	},
	{
		Name:        CmdInventory,
		Description: "Check your items that you've bought",
		Options: []*discordgo.ApplicationCommandOption{
			{
				Type:        discordgo.ApplicationCommandOptionUser,
				Name:        "user",
				Description: "Person whose inventory to show (optional)",
				Required:    false,
			},
		},
	},
	{
		Name:        CmdLeaderboard,
		Description: "Show top-10 users with diamonds",
	},
	{
		Name:        CmdRoulette,
		Description: "Try to win and not lose your money!",
		Options: []*discordgo.ApplicationCommandOption{
			{
				Type:        discordgo.ApplicationCommandOptionInteger,
				Name:        "amount",
				Description: "Amount of money to gamble",
				Required:    true,
			},
		},
	},
	{
		Name:        CmdTransfer,
		Description: "Transfer your money to another person",
		Options: []*discordgo.ApplicationCommandOption{
			{
				Type:        discordgo.ApplicationCommandOptionUser,
				Name:        "user",
				Description: "Person to whom to transfer",
				Required:    true,
			},
			{
				Type:        discordgo.ApplicationCommandOptionInteger,
				Name:        "amount",
				Description: "Amount of money to transfer",
				Required:    true,
			},
		},
	},
	{
		Name:        CmdWork,
		Description: "Work and get paid",
	},
	{
		Name:        CmdSteal,
		Description: "Steal money from someone; you can fail though, be careful",
		Options: []*discordgo.ApplicationCommandOption{
			{
				Type:        discordgo.ApplicationCommandOptionUser,
				Name:        "user",
				Description: "Person to steal money from",
				Required:    true,
			},
		},
	},
	{
		Name:        CmdBuy,
		Description: "Buy an item",
		Options: []*discordgo.ApplicationCommandOption{
			{
				Type:        discordgo.ApplicationCommandOptionString,
				Name:        "item",
				Description: "Item to buy",
				Required:    true,
			},
			{
				Type:        discordgo.ApplicationCommandOptionInteger,
				Name:        "amount",
				Description: "Amount of items to buy (optional)",
				Required:    false,
			},
		},
	},
	{
		Name:        CmdSell,
		Description: "Sell an item",
		Options: []*discordgo.ApplicationCommandOption{
			{
				Type:        discordgo.ApplicationCommandOptionString,
				Name:        "item",
				Description: "Item to sell",
				Required:    true,
			},
			{
				Type:        discordgo.ApplicationCommandOptionInteger,
				Name:        "amount",
				Description: "Amount of items to sell (optional)",
				Required:    false,
			},
		},
	},
	{
		Name:        CmdGive,
		Description: "Give an item to someone",
		Options: []*discordgo.ApplicationCommandOption{
			{
				Type:        discordgo.ApplicationCommandOptionUser,
				Name:        "user",
				Description: "Person to give item to",
				Required:    true,
			},
			{
				Type:        discordgo.ApplicationCommandOptionString,
				Name:        "item",
				Description: "Item to give",
				Required:    true,
			},
			{
				Type:        discordgo.ApplicationCommandOptionInteger,
				Name:        "amount",
				Description: "Amount of items to give (optional)",
				Required:    false,
			},
		},
	},
	{
		Name:        CmdEat,
		Description: "Eat something from your inventory",
		Options: []*discordgo.ApplicationCommandOption{
			{
				Type:        discordgo.ApplicationCommandOptionString,
				Name:        "item",
				Description: "Item to eat",
				Required:    true,
			},
		},
	},
	{
		Name:        CmdCooldown,
		Description: "Check your cooldowns",
		Options: []*discordgo.ApplicationCommandOption{
			{
				Type:        discordgo.ApplicationCommandOptionUser,
				Name:        "user",
				Description: "Person to check (optional)",
				Required:    false,
			},
		},
	},
	{
		Name:        CmdHigh,
		Description: "Check if you are high and the remaining time if you are",
		Options: []*discordgo.ApplicationCommandOption{
			{
				Type:        discordgo.ApplicationCommandOptionUser,
				Name:        "user",
				Description: "Person to check (optional)",
				Required:    false,
			},
		},
	},
}
