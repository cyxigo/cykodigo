package bot

import "github.com/bwmarrin/discordgo"

const (
	cmdHelp      = "help"
	cmdMe        = "me"
	cmdMeow      = "meow"
	cmdMeowat    = "meowat"
	cmdBark      = "bark"
	cmdBarkat    = "barkat"
	cmdDoflip    = "doflip"
	cmdExplode   = "explode"
	cmdSpin      = "spin"
	cmdCat       = "cat"
	cmdCart      = "cart"
	cmdRoulette  = "roulette"
	cmdAssault   = "assault"
	cmdWork      = "work"
	cmdBalance   = "balance"
	cmdTransfer  = "transfer"
	cmdSteal     = "steal"
	cmdShop      = "shop"
	cmdBuy       = "buy"
	cmdInventory = "inventory"
	cmdEat       = "eat"

	// msg means that this is message content command
	// same for handlers: handleMsg
	cmdMsgMeow         = "meow"
	cmdMsgCrazy        = "crazy"
	cmdMsgExplodeBalls = "cykodigo explode balls"
	cmdMsgGlamptastic  = "glamptastic"
)

// slash commands, looooooooooong list of them
var cmds = []*discordgo.ApplicationCommand{
	{
		Name:        cmdHelp,
		Description: "This... will not help actually",
	},
	{
		Name:        cmdMe,
		Description: "Send picture of me! Nah, not you of course",
	},
	{
		Name:        cmdMeow,
		Description: "He will meow",
	},
	{
		Name:        cmdMeowat,
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
		Name:        cmdBark,
		Description: "He will... or won't bark. I don't know he's a cat",
	},
	{
		Name:        cmdBarkat,
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
		Name:        cmdDoflip,
		Description: "He will do a flip",
	},
	{
		Name:        cmdExplode,
		Description: "He will explo- WHAT??? DON'T DO THAT!!!",
	},
	{
		Name:        cmdSpin,
		Description: "He will spin! Wooooooo",
	},
	{
		Name:        cmdCat,
		Description: "Cat!",
	},
	{
		Name:        cmdCart,
		Description: "Cart! (totally NOT copied from Cat Bot hehe)",
	},
	{
		Name:        cmdRoulette,
		Description: "Why don't we play a little game?",
	},
	{
		Name:        cmdAssault,
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
		Name:        cmdWork,
		Description: "Work and get paid! Money!1!11!!",
	},
	{
		Name:        cmdBalance,
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
		Name:        cmdTransfer,
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
		Name:        cmdSteal,
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
		Name:        cmdShop,
		Description: "Shop!!! Buy things, lose your money!1!11!!",
	},
	{
		Name:        cmdBuy,
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
		Name:        cmdInventory,
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
		Name:        cmdEat,
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
