package data

import "github.com/bwmarrin/discordgo"

const (
	CmdHelp        = "help"
	CmdMe          = "me"
	CmdMeow        = "meow"
	CmdMeowat      = "meowat"
	CmdBark        = "bark"
	CmdBarkat      = "barkat"
	CmdDoflip      = "doflip"
	CmdExplode     = "explode"
	CmdSpin        = "spin"
	CmdCat         = "cat"
	CmdCart        = "cart"
	CmdRoulette    = "roulette"
	CmdAssault     = "assault"
	CmdWork        = "work"
	CmdBalance     = "balance"
	CmdTransfer    = "transfer"
	CmdSteal       = "steal"
	CmdShop        = "shop"
	CmdBuy         = "buy"
	CmdInventory   = "inventory"
	CmdLeaderboard = "leaderboard"
	CmdEat         = "eat"

	// msg means that this is message content command
	// same for handlers: handleMsg
	CmdMsgMeow         = "meow"
	CmdMsgCrazy        = "crazy"
	CmdMsgExplodeBalls = "cykodigo explode balls"
	CmdMsgGlamptastic  = "glamptastic"
	CmdMsgNature       = "nature"
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
		Description: "Meow at someone!",
		Options: []*discordgo.ApplicationCommandOption{
			{
				Type:        discordgo.ApplicationCommandOptionUser,
				Name:        "member",
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
		Description: "Bark at someone!", Options: []*discordgo.ApplicationCommandOption{
			{
				Type:        discordgo.ApplicationCommandOptionUser,
				Name:        "member",
				Description: "Person to bark at",
				Required:    true,
			},
		},
	},
	{
		Name:        CmdDoflip,
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
		Name:        CmdCat,
		Description: "Cat!",
	},
	{
		Name:        CmdCart,
		Description: "Cart! (totally NOT copied from Cat Bot hehe)",
	},
	{
		Name:        CmdRoulette,
		Description: "Why don't we play a little game?",
	},
	{
		Name:        CmdAssault,
		Description: "Try to assault someone... shhh...",
		Options: []*discordgo.ApplicationCommandOption{
			{
				Type:        discordgo.ApplicationCommandOptionUser,
				Name:        "member",
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
		Name:        CmdWork,
		Description: "Work and get paid! Money!1!11!!",
	},
	{
		Name:        CmdBalance,
		Description: "Balance! Check how much money you've got from hard work",
		Options: []*discordgo.ApplicationCommandOption{
			{
				Type:        discordgo.ApplicationCommandOptionUser,
				Name:        "member",
				Description: "Person whose balance to show (optional)",
				Required:    false,
			},
		},
	},
	{
		Name:        CmdTransfer,
		Description: "Transfer your money to another person",
		Options: []*discordgo.ApplicationCommandOption{
			{
				Type:        discordgo.ApplicationCommandOptionUser,
				Name:        "member",
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
		Name:        CmdSteal,
		Description: "Steal money from someone! You can fail though, be careful.",
		Options: []*discordgo.ApplicationCommandOption{
			{
				Type:        discordgo.ApplicationCommandOptionUser,
				Name:        "member",
				Description: "Person to steal money from",
				Required:    true,
			},
		},
	},
	{
		Name:        CmdShop,
		Description: "Shop!!! Buy things, lose your money!1!11!!",
	},
	{
		Name:        CmdBuy,
		Description: "Buy item!1!11!!", Options: []*discordgo.ApplicationCommandOption{
			{
				Type:        discordgo.ApplicationCommandOptionString,
				Name:        "item",
				Description: "Item to buy",
				Required:    true,
			},
		},
	},
	{
		Name:        CmdInventory,
		Description: "Inventory! Check your items that you've bought",
		Options: []*discordgo.ApplicationCommandOption{
			{
				Type:        discordgo.ApplicationCommandOptionUser,
				Name:        "member",
				Description: "Person whose inventory to show (optional)",
				Required:    false,
			},
		},
	},
	{
		Name:        CmdLeaderboard,
		Description: "Leaderboard! Compete!!! Diamonds!1!11!!",
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
}
